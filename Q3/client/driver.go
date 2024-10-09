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
    "sync"
    "time"

    "go.etcd.io/etcd/client/v3"
    pb "github.com/Mitanshk01/DS_HW4/Q3/protofiles"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "github.com/google/uuid"
)

var (
    rideInProgress bool
    driverID       string
	activeRideID   string
	currentClient  pb.DriverServiceClient
    rideQueueMutex sync.Mutex
)

func discoverServers(etcdClient *clientv3.Client) []string {
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

func isPortAvailable(port string) bool {
    ln, err := net.Listen("tcp", ":"+port)
    if err != nil {
        return false
    }
    ln.Close()
    return true
}

func handleRideRequests(client pb.DriverServiceClient, wg *sync.WaitGroup) {
    defer wg.Done()

    req := &pb.AssignmentRequest{
        DriverId: driverID,
    }

    stream, err := client.GetRideRequest(context.Background(), req)
    if err != nil {
        log.Fatalf("Error while connecting to ride request stream: %v\n", err)
    }

    for {
        res, err := stream.Recv()
        if err != nil {
            log.Printf("Error receiving ride request from stream: %v\n", err)
            return
        }
		
		rideQueueMutex.Lock()
		
		if !rideInProgress {
			fmt.Printf("Offering Ride ID: %s to you.\n", res.RideId)
            fmt.Printf("Pickup Location: %s\n", res.PickupLocation)
            fmt.Printf("Destination: %s\n", res.Destination)
			
			for {
				fmt.Println("Looping")
				fmt.Println("Type 'accept' to Accept the Ride, or 'reject' to Reject the Ride:")
	
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				choice := strings.ToLower(strings.TrimSpace(scanner.Text()))
	
				if choice == "accept" {
					acceptRide(client, res.RideId)
					break
				} else if choice == "reject" {
					rejectRide(client, res.RideId)
					break
				} else {
					fmt.Println("Invalid choice. Please type 'accept' or 'reject'.")
				}
			}	
		}

		rideQueueMutex.Unlock()

        time.Sleep(1 * time.Second)
    }
}

func connectToAllServers(serverAddresses []string, tlsCredentials credentials.TransportCredentials, wg *sync.WaitGroup) {
	for _, serverAddress := range serverAddresses {
		fullServerAddress := fmt.Sprintf("127.0.0.1:%s", serverAddress)
		conn, err := grpc.Dial(fullServerAddress, grpc.WithTransportCredentials(tlsCredentials))
		if err != nil {
			log.Printf("Failed to connect to server %s: %v", serverAddress, err)
			continue
		}

		client := pb.NewDriverServiceClient(conn)

		wg.Add(1)
		go handleRideRequests(client, wg)
	}
}

func driverCLI() {
    for {
		if rideInProgress {
			fmt.Println("\nDriver Menu:");
			scanner := bufio.NewScanner(os.Stdin)
        	fmt.Println("Type 'complete' to complete the current Ride (if any)");
			fmt.Println("Type 'exit' to log out and exit");

			fmt.Print("\nEnter your choice: ");
        	scanner.Scan();
			choice := strings.ToLower(strings.TrimSpace(scanner.Text()));
			switch choice {
			case "complete":
				if rideInProgress {
					completeRide()
				} else {
					fmt.Println("Error: No active ride to complete.")
				}
			case "exit":
				if rideInProgress {
					fmt.Println("Error: You cannot exit while in an ongoing ride. Complete the ride first.")
				} else {
					fmt.Println("Logging out and exiting the system.")
					return
				}
			default:
				fmt.Println("Invalid choice. Please type 'complete' or 'exit'.")
			}
		}
    }
}

func acceptRide(client pb.DriverServiceClient, rideID string) {
    req := &pb.AcceptRideRequest{
        DriverId: driverID,
        RideId:   rideID,
    }

	fmt.Println("Accepting Ride")
    res, err := client.AcceptRide(context.Background(), req)
    if err != nil {
        log.Fatalf("Failed to accept ride: %v", err)
    }

    if res.Status == "Accepted" {
        rideInProgress = true
		activeRideID = rideID
		currentClient = client
        fmt.Printf("Ride ID: %s accepted. Ride is now in progress.\n", rideID)
    } else {
        fmt.Printf("Ride ID: %s could not be accepted: %s\n", rideID, res.Status)
    }
}

func rejectRide(client pb.DriverServiceClient, rideID string) {
    req := &pb.RejectRideRequest{
        DriverId: driverID,
        RideId:   rideID,
    }

    res, err := client.RejectRide(context.Background(), req)
    if err != nil {
        log.Fatalf("Failed to reject ride: %v", err)
    }

    fmt.Printf("Ride ID: %s rejected: %s\n", rideID, res.Status)
}

func completeRide() {
    if activeRideID == "" {
        fmt.Println("Error: No active ride to complete.")
        return
    }

    req := &pb.RideCompletionRequest{
        DriverId: driverID,
        RideId:   activeRideID,
    }

    res, err := currentClient.CompleteRide(context.Background(), req)
    if err != nil {
        log.Fatalf("Failed to complete ride: %v", err)
    }

    if res.Status == "Completed" {
        rideInProgress = false
		activeRideID = ""
		currentClient = nil
        fmt.Printf("Ride ID: %s completed.\n", activeRideID)
    } else {
        fmt.Printf("Ride ID: %s could not be completed: %s\n", activeRideID, res.Status)
    }
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	clientCert, err := tls.LoadX509KeyPair("../certificate/driver-cert.pem", "../certificate/driver-key.pem")
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

func main() {
	port := flag.String("port", "", "Port to connect to the server (e.g., 50051)")
	etcdURL := flag.String("etcd", "localhost:2379", "Etcd server URL")
	flag.Parse()

	if *port == "" {
		log.Fatal("Error: Port number is required. Please provide it using the --port flag.")
	}

	if !isPortAvailable(*port) {
		log.Fatalf("Error: Port %s is already in use or unavailable.", *port)
	}

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{*etcdURL},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Error creating etcd client: %v", err)
	}
	defer etcdClient.Close()

	serverAddresses := discoverServers(etcdClient)
	if len(serverAddresses) == 0 {
		log.Fatalf("No servers found in etcd.")
	}

	fmt.Printf("Discovered servers: %v\n", serverAddresses)

	driverID = uuid.New().String()

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("Cannot load TLS credentials: ", err)
	}

	var wg sync.WaitGroup

	connectToAllServers(serverAddresses, tlsCredentials, &wg)

	driverCLI()

	wg.Wait()
}
