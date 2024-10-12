package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"go.etcd.io/etcd/client/v3"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
	pb "github.com/Mitanshk01/DS_HW4/Q3/protofiles"
)

var rideInProgress bool = false
var currentRideID string = ""
var riderID string

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	clientCert, err := tls.LoadX509KeyPair("./certificate/rider-cert.pem", "./certificate/rider-key.pem")
	if err != nil {
		log.Fatalf("Failed to load client certificate: %v", err)
	}

	caCert, err := ioutil.ReadFile("./certificate/ca-cert.pem")
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
		ServerName:   "localhost", 
	}

	return credentials.NewTLS(tlsConfig), nil
}

func discoverServers() []string {
	etcdAddress := flag.String("etcd", "http://localhost:2379", "etcd server address")
	flag.Parse()

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{*etcdAddress},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer etcdClient.Close()

	resp, err := etcdClient.Get(context.Background(), "/ride-sharing/servers", clientv3.WithPrefix())
	if err != nil {
		log.Fatalf("Error retrieving servers from etcd: %v", err)
	}

	var servers []string
	for _, kv := range resp.Kvs {
		servers = append(servers, fmt.Sprintf("localhost:%s", string(kv.Value)))
	}	

	return servers
}

func riderCLI(client pb.RiderServiceClient) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nRider Menu:")
		fmt.Println("Type 'request' to Request a Ride")
		fmt.Println("Type 'status' to Check Ride Status")
		fmt.Println("Type 'exit' to Log Out and Exit")

		fmt.Print("\nEnter your choice: ")
		scanner.Scan()
		choice := strings.ToLower(strings.TrimSpace(scanner.Text()))

		switch choice {
		case "request":
			if rideInProgress {
				fmt.Println("Error: You already have an active ride. Complete the ride before requesting a new one.")
			} else {
				requestRide(client)
			}
		case "status":
			if rideInProgress {
				checkRideStatus(client)
			} else {
				fmt.Println("Error: No active ride. Request a ride first.")
			}
		case "exit":
			if rideInProgress {
				fmt.Println("Error: You cannot exit while in an ongoing ride. Complete the ride first.")
			} else {
				fmt.Println("Logging out and exiting the system.")
				return
			}
		default:
			fmt.Println("Invalid choice. Please type 'request', 'status', or 'exit'.")
		}
	}
}

func requestRide(client pb.RiderServiceClient) {
	fmt.Print("Enter Pickup Location: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	pickup := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter Destination Location: ")
	scanner.Scan()
	destination := strings.TrimSpace(scanner.Text())

	req := &pb.RideRequest{
		RiderId:         riderID,
		PickupLocation:  pickup,
		Destination:     destination,
	}

	fmt.Println("Ride Requested")
	res, err := client.RequestRide(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to request ride: %v", err)
	}

	if res.Status == "Ongoing" {
		rideInProgress = true
		currentRideID = res.RideId
		fmt.Printf("Ride requested. Ride ID: %s. Driver ID: %s assigned.\n", res.RideId, res.DriverId)
	} else {
		fmt.Printf("Ride request failed: %s\n", res.Status)
	}
}

func checkRideStatus(client pb.RiderServiceClient) {
	req := &pb.RideStatusRequest{
		RideId: currentRideID,
	}

	res, err := client.GetRideStatus(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to get ride status: %v", err)
	}

	fmt.Printf("Ride ID: %s, Status: %s\n", currentRideID, res.Status)

	if res.Status == "Completed" {
		rideInProgress = false
		fmt.Println("The ride has been completed by the driver.")
		currentRideID = ""
	}
}

func main() {
	policy := flag.String("policy", "round_robin", "Load balancing policy: 'pick_first' or 'round_robin'")
	flag.Parse()

	servers := discoverServers()
	fmt.Printf("Discovered servers: %v\n", servers)

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("Cannot load TLS credentials: ", err)
	}

	r := manual.NewBuilderWithScheme("myservice")
	var addresses []resolver.Address
	for _, server := range servers {
		addresses = append(addresses, resolver.Address{Addr: server})
	}
	r.InitialState(resolver.State{Addresses: addresses})

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(tlsCredentials),
		grpc.WithResolvers(r),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{
			"loadBalancingPolicy": "%s",
			"methodConfig": [{
				"name": [{"service": ""}],
				"retryPolicy": {
					"MaxAttempts": 5,
					"InitialBackoff": "0.1s",
					"MaxBackoff": "1s",
					"BackoffMultiplier": 2.0,
					"RetryableStatusCodes": ["UNAVAILABLE"]
				}
			}] 
		}`, *policy)),
	}

	conn, err := grpc.Dial("myservice:///", dialOpts...)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewRiderServiceClient(conn)

	riderID = uuid.New().String()
	fmt.Printf("Your unique Rider ID: %s\n", riderID)

	riderCLI(client)
}