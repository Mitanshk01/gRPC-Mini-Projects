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
	"strings"
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
	mu                   sync.Mutex
	ch                   sync.Mutex
	driverStreams        map[string]pb.DriverService_GetRideRequestServer
	maxRetries           int
	etcdClient           *clientv3.Client
}

func NewRideSharingServer(etcdClient *clientv3.Client) *rideSharingServer {
	return &rideSharingServer{
		driverStreams:        make(map[string]pb.DriverService_GetRideRequestServer),
		maxRetries:           3,
		etcdClient:           etcdClient,
	}
}

func createAssignmentKey(rideID string, driverID string) string {
	return rideID + "_" + driverID
}

func (s *rideSharingServer) setRideStatus(rideID, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := s.etcdClient.Put(ctx, fmt.Sprintf("/ride-sharing/ride-status/%s", rideID), status)
	return err
}

func (s *rideSharingServer) getRideStatus(rideID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := s.etcdClient.Get(ctx, fmt.Sprintf("/ride-sharing/ride-status/%s", rideID))
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("ride not found")
	}
	return string(resp.Kvs[0].Value), nil
}

func (s *rideSharingServer) setRideAssignmentStatus(rideID, driverID, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	key := fmt.Sprintf("/ride-sharing/ride-assignment/%s_%s", rideID, driverID)
	_, err := s.etcdClient.Put(ctx, key, status)
	return err
}

func (s *rideSharingServer) getRideAssignmentStatus(rideID, driverID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	key := fmt.Sprintf("/ride-sharing/ride-assignment/%s_%s", rideID, driverID)
	resp, err := s.etcdClient.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("assignment not found")
	}
	return string(resp.Kvs[0].Value), nil
}

func (s *rideSharingServer) setDriverStatus(driverID string, available bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status := "unavailable"
	if available {
		status = "available"
	}

	_, err := s.etcdClient.Put(ctx, "/ride-sharing/drivers/"+driverID, status)
	return err
}

func (s *rideSharingServer) getDriverStatus(driverID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := s.etcdClient.Get(ctx, "/ride-sharing/drivers/"+driverID)
	if err != nil {
		return false, err
	}

	if len(resp.Kvs) == 0 {
		return false, nil
	}

	return string(resp.Kvs[0].Value) == "available", nil
}

func (s *rideSharingServer) getAvailableDrivers() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := s.etcdClient.Get(ctx, "/ride-sharing/drivers/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var availableDrivers []string
	for _, kv := range resp.Kvs {
		if string(kv.Value) == "available" {
			driverID := strings.TrimPrefix(string(kv.Key), "/ride-sharing/drivers/")
			availableDrivers = append(availableDrivers, driverID)
		}
	}

	return availableDrivers, nil
}

func (s *rideSharingServer) RequestRide(ctx context.Context, req *pb.RideRequest) (*pb.RideResponse, error) {
	rideID := uuid.New().String()

	fmt.Printf("Request Arrived %s\n", rideID)
	driverID, err := s.assignDriver(rideID, req.RiderId, req.PickupLocation, req.Destination)
	if err != nil {
		fmt.Println(err)
		if err := s.setRideStatus(rideID, "Cancelled"); err != nil {
			log.Printf("Failed to set ride status: %v", err)
		}
		return &pb.RideResponse{Status: "No Drivers Available"}, nil
	}

	if err := s.setRideStatus(rideID, "Ongoing"); err != nil {
		log.Printf("Failed to set ride status: %v", err)
	}

	return &pb.RideResponse{RideId: rideID, DriverId: driverID, Status: "Ongoing"}, nil
}

func (s *rideSharingServer) GetRideRequest(req *pb.AssignmentRequest, stream pb.DriverService_GetRideRequestServer) error {
	fmt.Printf("Driver arrived with id: %s\n", req.DriverId)
	
	if err := s.setDriverStatus(req.DriverId, true); err != nil {
		log.Printf("Failed to set driver status: %v", err)
		return err
	}

	s.mu.Lock()
	s.driverStreams[req.DriverId] = stream
	s.mu.Unlock()

	<-stream.Context().Done()

	fmt.Println("Driver disconnected:", req.DriverId)
	s.mu.Lock()
	delete(s.driverStreams, req.DriverId)
	s.mu.Unlock()

	if err := s.setDriverStatus(req.DriverId, false); err != nil {
		log.Printf("Failed to update driver status on disconnect: %v", err)
	}

	return nil
}

func (s *rideSharingServer) assignDriver(rideID, riderID, pickupLocation, destination string) (string, error) {
	for attempt := 0; attempt < s.maxRetries; attempt++ {
		availableDrivers, err := s.getAvailableDrivers()
		if err != nil {
			log.Printf("Failed to get available drivers: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		for _, driverID := range availableDrivers {
			if err := s.setRideAssignmentStatus(rideID, driverID, "Pending"); err != nil {
				log.Printf("Failed to set ride assignment status: %v", err)
				continue
			}

			req := &pb.AssignmentResponse{
				RideId:         rideID,
				RiderId:        riderID,
				PickupLocation: pickupLocation,
				Destination:    destination,
			}

			driverStream, exists := s.driverStreams[driverID]
			if !exists {
				log.Printf("Driver %s not connected", driverID)
				if err := s.setRideAssignmentStatus(rideID, driverID, "Rejected"); err != nil {
					log.Printf("Failed to update ride assignment status: %v", err)
				}
				continue
			}

			go func(driverID string, stream pb.DriverService_GetRideRequestServer) {
				if err := stream.Send(req); err != nil {
					log.Printf("Failed to send assignment response to driver %s: %v", driverID, err)
					if err := s.setRideAssignmentStatus(rideID, driverID, "Rejected"); err != nil {
						log.Printf("Failed to update ride assignment status: %v", err)
					}
				}
			}(driverID, driverStream)

			time.Sleep(10 * time.Second)

			status, err := s.getRideAssignmentStatus(rideID, driverID)
			if err != nil {
				log.Printf("Failed to get ride assignment status: %v", err)
				continue
			}

			if status == "Accepted" {
				log.Printf("Ride %s accepted by driver %s\n", rideID, driverID)
				return driverID, nil
			} else {
				log.Printf("Ride %s not accepted by driver %s\n", rideID, driverID)
				if err := s.setRideAssignmentStatus(rideID, driverID, "Rejected"); err != nil {
					log.Printf("Failed to update ride assignment status: %v", err)
				}
			}
		}

		log.Printf("No drivers available currently, trying again...\n")
		time.Sleep(3 * time.Second)
	}

	return "", fmt.Errorf("no available drivers after multiple attempts")
}

func (s *rideSharingServer) AcceptRide(ctx context.Context, req *pb.AcceptRideRequest) (*pb.AcceptRideResponse, error) {
	status, err := s.getRideAssignmentStatus(req.RideId, req.DriverId)
	if err != nil {
		return &pb.AcceptRideResponse{Status: "Error getting ride status"}, err
	}

	if status == "Pending" {
		if err := s.setRideAssignmentStatus(req.RideId, req.DriverId, "Accepted"); err != nil {
			return &pb.AcceptRideResponse{Status: "Error updating ride status"}, err
		}
		
		if err := s.setDriverStatus(req.DriverId, false); err != nil {
			log.Printf("Failed to set driver %s as unavailable: %v", req.DriverId, err)
			return &pb.AcceptRideResponse{Status: "Error updating driver status"}, err
		}
		
		log.Printf("Ride %s accepted by driver %s", req.RideId, req.DriverId)
		return &pb.AcceptRideResponse{Status: "Accepted"}, nil
	}

	return &pb.AcceptRideResponse{Status: "Ride not found or already assigned"}, nil
}

func (s *rideSharingServer) RejectRide(ctx context.Context, req *pb.RejectRideRequest) (*pb.RejectRideResponse, error) {
    kvc := clientv3.NewKV(s.etcdClient)

    assignmentKey := fmt.Sprintf("/ride-sharing/ride-assignment/%s_%s", req.RideId, req.DriverId)
    driverKey := fmt.Sprintf("/ride-sharing/drivers/%s", req.DriverId)

    txn := kvc.Txn(ctx).If(
        clientv3.Compare(clientv3.Value(assignmentKey), "=", "Pending"),
    ).Then(
        clientv3.OpPut(assignmentKey, "Rejected"),
        clientv3.OpPut(driverKey, "available"),
    )

    txnResp, err := txn.Commit()
    if err != nil {
        log.Printf("Error in RejectRide transaction: %v", err)
        return &pb.RejectRideResponse{Status: "Error rejecting ride"}, err
    }

    if !txnResp.Succeeded {
        return &pb.RejectRideResponse{Status: "Ride not found or already processed"}, nil
    }

    log.Printf("Ride %s rejected by driver %s", req.RideId, req.DriverId)
    return &pb.RejectRideResponse{Status: "Rejected"}, nil
}


func (s *rideSharingServer) GetRideStatus(ctx context.Context, req *pb.RideStatusRequest) (*pb.RideStatusResponse, error) {
	status, err := s.getRideStatus(req.RideId)
	if err != nil {
		return &pb.RideStatusResponse{Status: "Ride not found"}, nil
	}

	return &pb.RideStatusResponse{RideId: req.RideId, Status: status}, nil
}

func (s *rideSharingServer) UpdateRideStatus(ctx context.Context, req *pb.UpdateRideStatusRequest) (*pb.UpdateRideStatusResponse, error) {
	if err := s.setRideStatus(req.RideId, req.Status); err != nil {
		return &pb.UpdateRideStatusResponse{Status: "Error updating ride status"}, err
	}
	return &pb.UpdateRideStatusResponse{Status: "Status updated to " + req.Status}, nil
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	serverCert, err := tls.LoadX509KeyPair("./certificate/server-cert.pem", "./certificate/server-key.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %v", err)
	}

	caCert, err := ioutil.ReadFile("./certificate/ca-cert.pem")
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

				if time.Now().After(cert.NotAfter) {
                    return nil, status.Error(codes.Unauthenticated, "certificate expired")
                }

                clientType := cert.Subject.CommonName

                method := info.FullMethod
                if strings.Contains(method, "DriverService") {
                    if clientType != "driver.rideshare.com" {
                        return nil, status.Errorf(codes.PermissionDenied, "unauthorized access for client type: %s", clientType)
                    }
                } else if strings.Contains(method, "RiderService") {
                    if clientType != "rider.rideshare.com" {
                        return nil, status.Errorf(codes.PermissionDenied, "unauthorized access for client type: %s", clientType)
                    }
                } else {
                    return nil, status.Error(codes.PermissionDenied, "invalid service")
                }

                return handler(ctx, req)
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
	lease, err := etcdClient.Grant(context.Background(), 60)

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
	if err := s.setRideStatus(req.RideId, "Completed"); err != nil {
		return &pb.RideCompletionResponse{Status: "Error completing ride"}, err
	}

	if err := s.setDriverStatus(req.DriverId, true); err != nil {
		log.Printf("Failed to set driver status to available: %v", err)
	}

	if err := s.setRideAssignmentStatus(req.RideId, req.DriverId, "Completed"); err != nil {
		log.Printf("Failed to update ride assignment status: %v", err)
	}

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

	// healthServer := health.NewServer()
	// grpc_health_v1.RegisterHealthServer(server, healthServer)

	// healthServer.SetServingStatus("RiderService", grpc_health_v1.HealthCheckResponse_SERVING)

	rideServer := NewRideSharingServer(etcdClient)
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
