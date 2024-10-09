package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"time"

	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "github.com/Mitanshk01/DS_HW4/Q3/protofiles"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"github.com/google/uuid"
)

type rideSharingServer struct {
	pb.UnimplementedRiderServiceServer
	pb.UnimplementedDriverServiceServer
	mu              sync.Mutex
	ch              sync.Mutex
	driverStatus    map[string]bool
	driverStreams   map[string]pb.DriverService_GetRideRequestServer
	rideStatusStore map[string]string
	rideAssignmentStatus map[string]string
	maxRetries      int
}

func NewRideSharingServer() *rideSharingServer {
	return &rideSharingServer{
		driverStatus:    make(map[string]bool), // Initialize with no drivers
		driverStreams:   make(map[string]pb.DriverService_GetRideRequestServer), // Store streams
		rideStatusStore: make(map[string]string),
		rideAssignmentStatus: make(map[string]string),
		maxRetries:      3,
	}
}

func createAssignmentKey(rideID string, driverID string) string {
	return rideID + "_" + driverID
}

func (s *rideSharingServer) RequestRide(ctx context.Context, req *pb.RideRequest) (*pb.RideResponse, error) {
	rideID := uuid.New().String()

	fmt.Println("Request Arrived %s", rideID)
	driverID, err := s.assignDriver(rideID, req.RiderId, req.PickupLocation, req.Destination)
	if err != nil {
		fmt.Println(err)
		return &pb.RideResponse{Status: "No drivers available"}, nil
	}

	s.mu.Lock()
	s.rideStatusStore[rideID] = "Assigned"
	s.mu.Unlock()

	return &pb.RideResponse{RideId: rideID, DriverId: driverID, Status: "Assigned"}, nil
}



func (s *rideSharingServer) GetRideRequest(req *pb.AssignmentRequest, stream pb.DriverService_GetRideRequestServer) error {
	s.mu.Lock()
	fmt.Println("Driver arrived with id: %s", req.DriverId)
	s.driverStatus[req.DriverId] = true
	s.driverStreams[req.DriverId] = stream
	s.mu.Unlock()

	<-stream.Context().Done()

	fmt.Println("Deleting driver")
	s.mu.Lock()
	delete(s.driverStreams, req.DriverId)
	s.driverStatus[req.DriverId] = false
	s.mu.Unlock()

	return nil
}


func (s *rideSharingServer) updateAssignmentStatus(assignmentKey, status string) {
    s.ch.Lock()
    s.rideAssignmentStatus[assignmentKey] = status
    s.ch.Unlock()
}


func (s *rideSharingServer) assignDriver(rideID, riderID, pickupLocation, destination string) (string, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

	fmt.Println("Requesting for %s %s", rideID, riderID)
    for attempt := 0; attempt < s.maxRetries; attempt++ {
        for driverID, available := range s.driverStatus {
            if available {
                assignmentKey := createAssignmentKey(rideID, driverID)

                s.ch.Lock()
                s.rideAssignmentStatus[assignmentKey] = "Pending" // Initially mark as pending
                s.ch.Unlock()

                s.driverStatus[driverID] = false // Mark driver as unavailable

                req := &pb.AssignmentResponse{
                    RideId:         rideID,
                    RiderId:        riderID,
                    PickupLocation: pickupLocation,
                    Destination:    destination,
                }

                driverStream, exists := s.driverStreams[driverID]
                if !exists {
                    log.Printf("Driver %s not connected", driverID)
                    s.updateAssignmentStatus(assignmentKey, "Rejected")
                    continue
                }

                go func(driverID string, stream pb.DriverService_GetRideRequestServer, assignmentKey string) {
                    if err := stream.Send(req); err != nil {
                        log.Printf("Failed to send assignment response to driver %s: %v", driverID, err)
                        s.updateAssignmentStatus(assignmentKey, "Rejected")
                    }
                }(driverID, driverStream, assignmentKey)

                time.Sleep(10 * time.Second)

                s.ch.Lock()
                res := s.rideAssignmentStatus[assignmentKey] // Check the current status
                s.ch.Unlock()

                if res == "Accepted" {
                    log.Printf("Ride %s accepted by driver %s\n", rideID, driverID)
                    return driverID, nil
                } else if res == "Rejected" {
					s.updateAssignmentStatus(assignmentKey, "Rejected")
                    log.Printf("Ride %s rejected by driver %s\n", rideID, driverID)
                    s.driverStatus[driverID] = true // Mark driver as available again
                } else {
                    log.Printf("Driver %s did not respond in time. Reassigning...\n", driverID)
                    s.updateAssignmentStatus(assignmentKey, "Rejected")
                    s.driverStatus[driverID] = true // Mark driver as available again
                }

            }
        }
    }

    return "", fmt.Errorf("no available drivers after multiple attempts")
}

func (s *rideSharingServer) AcceptRide(ctx context.Context, req *pb.AcceptRideRequest) (*pb.AcceptRideResponse, error) {
    assignmentKey := createAssignmentKey(req.RideId, req.DriverId)

    s.ch.Lock()
    defer s.ch.Unlock()

    if status, exists := s.rideAssignmentStatus[assignmentKey]; exists && status == "Pending" {
        s.rideAssignmentStatus[assignmentKey] = "Accepted"
        log.Println("Ride accepted by driver", req.DriverId)
        return &pb.AcceptRideResponse{Status: "Accepted"}, nil
    }

    return &pb.AcceptRideResponse{Status: "Ride not found or already assigned"}, nil
}

func (s *rideSharingServer) RejectRide(ctx context.Context, req *pb.RejectRideRequest) (*pb.RejectRideResponse, error) {
    assignmentKey := createAssignmentKey(req.RideId, req.DriverId)

    s.ch.Lock()
    defer s.ch.Unlock()

    if status, exists := s.rideAssignmentStatus[assignmentKey]; exists && status != "Pending" {
        s.rideAssignmentStatus[assignmentKey] = "Rejected"
        log.Println("Ride rejected by driver", req.DriverId)
        return &pb.RejectRideResponse{Status: "Rejected"}, nil
    }

    return &pb.RejectRideResponse{Status: "Ride not found or already assigned"}, nil
}


func (s *rideSharingServer) GetRideStatus(ctx context.Context, req *pb.RideStatusRequest) (*pb.RideStatusResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	status, exists := s.rideStatusStore[req.RideId]
	if !exists {
		return &pb.RideStatusResponse{Status: "Ride not found"}, nil
	}

	return &pb.RideStatusResponse{RideId: req.RideId, Status: status}, nil
}

func (s *rideSharingServer) UpdateRideStatus(ctx context.Context, req *pb.UpdateRideStatusRequest) (*pb.UpdateRideStatusResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.rideStatusStore[req.RideId]; !exists {
		return &pb.UpdateRideStatusResponse{Status: "Ride not found"}, nil
	}

	s.rideStatusStore[req.RideId] = req.Status
	return &pb.UpdateRideStatusResponse{Status: "Status updated to " + req.Status}, nil
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	serverCert, err := tls.LoadX509KeyPair("../certificate/server-cert.pem", "../certificate/server-key.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %v", err)
	}

	caCert, err := ioutil.ReadFile("../certificate/ca-cert.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	return credentials.NewTLS(tlsConfig), nil
}

func isPortAvailable(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

func authorizationInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if _, ok := metadata.FromIncomingContext(ctx); ok {
			peerInfo, _ := peer.FromContext(ctx)
			if tlsInfo, ok := peerInfo.AuthInfo.(credentials.TLSInfo); ok {
				cert := tlsInfo.State.PeerCertificates[0]

				clientType := cert.Subject.CommonName

				switch clientType {
				case "driver.rideshare.com":
					return handler(ctx, req)
				case "rider.rideshare.com":
					return handler(ctx, req)
				}

				return nil, status.Errorf(codes.PermissionDenied, "unauthorized client type: %s", clientType)
			}
		}

		return nil, status.Error(codes.Unauthenticated, "missing or invalid credentials")
	}
}

func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if clientRole := md["client-role"]; len(clientRole) > 0 {
				log.Printf("Time: %s; Method: %s; Client Role: %s; Request: %v",
					time.Now().Format(time.RFC3339), info.FullMethod, clientRole[0], req)
			} else {
				log.Printf("Time: %s; Method: %s; Client Role: Unknown; Request: %v",
					time.Now().Format(time.RFC3339), info.FullMethod, req)
			}
		} else {
			log.Printf("Time: %s; Method: %s; Client Role: Unknown; Request: %v",
				time.Now().Format(time.RFC3339), info.FullMethod, req)
		}

		resp, err := handler(ctx, req)

		log.Printf("Time: %s; Method: %s; Response: %v; Error: %v",
			time.Now().Format(time.RFC3339), info.FullMethod, resp, err)

		return resp, err
	}
}

func combinedInterceptor() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        authInterceptor := authorizationInterceptor()

        loggingHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
            return LoggingInterceptor()(ctx, req, info, handler)
        }

        return authInterceptor(ctx, req, info, loggingHandler)
    }
}

func registerServer(etcdClient *clientv3.Client, port string) error {
	lease, err := etcdClient.Grant(context.Background(), 60) // TTL of 60 seconds

	if err != nil {
		return fmt.Errorf("failed to create lease: %v", err)
	}

	_, err = etcdClient.Put(context.Background(), "/ride-sharing/servers/"+port, port, clientv3.WithLease(lease.ID))
	if err != nil {
		return fmt.Errorf("failed to register server with etcd: %v", err)
	}

	ch, err := etcdClient.KeepAlive(context.Background(), lease.ID)
	if err != nil {
		return fmt.Errorf("failed to keep etcd lease alive: %v", err)
	}

	go func() {
		for {
			ka := <-ch
			if ka == nil {
				log.Println("KeepAlive channel closed")
				return
			}
			log.Printf("Lease renewed for server at port: %s\n", port)
		}
	}()

	return nil
}

func (s *rideSharingServer) CompleteRide(ctx context.Context, req *pb.RideCompletionRequest) (*pb.RideCompletionResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if _, exists := s.rideStatusStore[req.RideId]; !exists {
        return &pb.RideCompletionResponse{Status: "Ride not found"}, nil
    }

    s.rideStatusStore[req.RideId] = "Completed"

    s.driverStatus[req.DriverId] = true

    assignmentKey := createAssignmentKey(req.RideId, req.DriverId)
    s.ch.Lock()
    delete(s.rideAssignmentStatus, assignmentKey)
    s.ch.Unlock()

    log.Printf("Ride %s completed by driver %s\n", req.RideId, req.DriverId)

    return &pb.RideCompletionResponse{Status: "Completed"}, nil
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
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer etcdClient.Close()

	if err := registerServer(etcdClient, *port); err != nil {
		log.Fatalf("Failed to register server with etcd: %v", err)
	}

	tlsCreds, err := loadTLSCredentials()
	if err != nil {
		log.Fatalf("Failed to load TLS credentials: %v", err)
	}

	server := grpc.NewServer(
		grpc.Creds(tlsCreds),
		grpc.UnaryInterceptor(combinedInterceptor()), 
	)

	rideServer := NewRideSharingServer()
	pb.RegisterRiderServiceServer(server, rideServer)
	pb.RegisterDriverServiceServer(server, rideServer)

	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", *port, err)
	}

	log.Printf("Starting secure gRPC server on port %s...\n", *port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
