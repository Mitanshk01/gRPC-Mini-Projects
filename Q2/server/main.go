package main

import (
    "flag"
    "log"
    "math"
    "net"
    "os"
    "bufio"
    "strconv"
    "strings"
    "sort"
    "fmt"

    "google.golang.org/grpc"
    knn "github.com/Mitanshk01/DS_HW4/Q2/protofiles"
)

type server struct {
    knn.UnimplementedKNNServiceServer
    dataset []knn.DataPoint
}

func calculateDistance(p1, p2 knn.DataPoint) float64 {
    var sum float64
    for i := 0; i < len(p1.Coordinates); i++ {
        diff := float64(p1.Coordinates[i] - p2.Coordinates[i])
        sum += diff * diff
    }
    return math.Sqrt(sum)
}

func (s *server) FindKNearestNeighbors(req *knn.KNNRequest, stream knn.KNNService_FindKNearestNeighborsServer) error {
    var neighbors []knn.Neighbor
    for _, point := range s.dataset {
        dist := calculateDistance(*req.QueryPoint, point)
        neighbors = append(neighbors, knn.Neighbor{
            Point:    &point,
            Distance: float32(dist),
        })
    }

    sort.Slice(neighbors, func(i, j int) bool {
        return neighbors[i].Distance < neighbors[j].Distance
    })

    for i := 0; i < int(req.K) && i < len(neighbors); i++ {
        neighbor := neighbors[i]
        if err := stream.Send(&neighbor); err != nil {
            return err
        }
    }

    return nil
}

func loadDatasetFromFile(filename string) ([]knn.DataPoint, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var dataset []knn.DataPoint
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        parts := strings.Fields(line)
        var coords []float32
        for _, val := range parts {
            coord, err := strconv.ParseFloat(val, 32)
            if err != nil {
                return nil, err
            }
            coords = append(coords, float32(coord))
        }
        dataset = append(dataset, knn.DataPoint{Coordinates: coords})
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return dataset, nil
}

func main() {
    port := flag.String("port", "", "Port number for the server")
    datasetFile := flag.String("dataset", "", "Path to the dataset file for the server")
    flag.Parse()

    if *port == "" || *datasetFile == "" {
        fmt.Println("Error: Both -port and -dataset flags are required.")
        flag.Usage()
        os.Exit(1)
    }

    datasetSubset, err := loadDatasetFromFile(*datasetFile)
    if err != nil {
        log.Fatalf("Failed to load dataset: %v", err)
    }

    lis, err := net.Listen("tcp", ":"+*port)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    knn.RegisterKNNServiceServer(grpcServer, &server{dataset: datasetSubset})
    log.Printf("Server listening on port %s", *port)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
