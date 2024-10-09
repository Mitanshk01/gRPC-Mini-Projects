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
	"net"
	"os"
	"strings"
	"time"

	"go.etcd.io/etcd/client/v3"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "github.com/Mitanshk01/DS_HW4/Q3/protofiles"
)

var rideInProgress bool = false
var currentRideID string = ""
var riderID string
var currentServer string

type LoadBalancer struct {
	servers []string
	index   int
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	clientCert, err := tls.LoadX509KeyPair("../certificate/rider-cert.pem", "../certificate/rider-key.pem")
	if err != nil {
		log.Fatalf("Failed to load client certificate: %v", err)
	}

	caCert, err := ioutil.ReadFile("../certificate/ca-cert.pem")
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	}

	return credentials.NewTLS(tlsConfig), nil
}

func discoverServers() []string {
	etcdAddress := flag.String("etcd", "http://localhost:2379", "etcd server address")
    flag.Parse()

	etcdClient, err := clientv3.New(clientv3.Config{
        Endpoints: []string{*etcdAddress},
        DialTimeout: 5 * time.Second,
    })

    var servers []string

    resp, err := etcdClient.Get(context.Background(), "/ride-sharing/servers", clientv3.WithPrefix())

    if err != nil {
        log.Fatalf("Error retrieving servers from etcd: %v\n", err)
    }

    for _, kv := range resp.Kvs {
        servers = append(servers, string(kv.Value))
    }

    return servers
}

func NewLoadBalancer(servers []string) *LoadBalancer {
	return &LoadBalancer{
		servers: servers,
		index:   0,
	}
}

func (lb *LoadBalancer) PickFirst() string {
	return lb.servers[0]
}

func (lb *LoadBalancer) RoundRobin() string {
	server := lb.servers[lb.index]
	lb.index = (lb.index + 1) % len(lb.servers)
	return server
}

func isPortAvailable(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

func riderCLI(client pb.RiderServiceClient, lb *LoadBalancer) {
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
				requestRide(client, lb)
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

func requestRide(client pb.RiderServiceClient, lb *LoadBalancer) {
	fmt.Print("Enter Pickup Location: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	pickup := strings.TrimSpace(scanner.Text())

	fmt.Print("Enter Destination Location: ")
	scanner.Scan()
	destination := strings.TrimSpace(scanner.Text())

	req := &pb.RideRequest{
		RiderId:     riderID,
		PickupLocation:      pickup,
		Destination: destination,
	}

	selectedServer := fmt.Sprintf("127.0.0.1:%s", lb.RoundRobin())

	tlsCredentials, err := loadTLSCredentials()

	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	conn, err := grpc.Dial(selectedServer, grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		log.Fatalf("Failed to connect to server %s: %v", selectedServer, err)
	}
	defer conn.Close()

	client = pb.NewRiderServiceClient(conn)

	fmt.Println("Ride Requested")
	res, err := client.RequestRide(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to request ride: %v", err)
	}

	fmt.Println("%s", res.Status)

	if res.Status == "Assigned" {
		rideInProgress = true
		currentRideID = res.RideId
		currentServer = selectedServer
		fmt.Printf("Ride requested. Ride ID: %s. Driver ID: %s assigned.\n", res.RideId, res.DriverId)
	} else {
		fmt.Printf("Ride request failed: %s\n", res.Status)
	}
}

func checkRideStatus(client pb.RiderServiceClient) {
	req := &pb.RideStatusRequest{
		RideId: currentRideID,
	}

	tlsCredentials, err := loadTLSCredentials()

	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	conn, err := grpc.Dial(currentServer, grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		log.Fatalf("Failed to connect to server %s: %v", currentServer, err)
	}
	defer conn.Close()

	client = pb.NewRiderServiceClient(conn)

	res, err := client.GetRideStatus(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to get ride status: %v", err)
	}

	fmt.Printf("Ride ID: %s, Status: %s\n", currentRideID, res.Status)

	if res.Status == "Completed" {
		rideInProgress = false
		fmt.Println("The ride has been completed by the driver.")
		currentRideID = ""
		currentServer = ""
	}
}

func main() {
	port := flag.String("port", "", "Port to connect to the server (e.g., 50051)")
	flag.Parse()

	if *port == "" {
		log.Fatal("Error: Port number is required. Please provide it using the --port flag.")
	}

	if !isPortAvailable(*port) {
		log.Fatalf("Error: Port %s is already in use or unavailable.", *port)
	}

	servers := discoverServers()

	fmt.Printf("Servers: %v\n", servers)

	lb := NewLoadBalancer(servers)

	serverAddress := fmt.Sprintf("127.0.0.1:%s", *port)
	
	tlsCredentials, err := loadTLSCredentials()

	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewRiderServiceClient(conn)

	riderID = uuid.New().String()
	fmt.Printf("Your unique Rider ID: %s\n", riderID)

	riderCLI(client, lb)
}
