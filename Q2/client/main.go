package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "sort"
    "sync"
    "time"

    "google.golang.org/grpc"
    knn "github.com/Mitanshk01/DS_HW4/Q2/protofiles"
)

const (
    numServers = 5
    basePort   = 50051
)

type Neighbor struct {
    Point    *knn.DataPoint
    Distance float32
}

type ByDistance []Neighbor

func (a ByDistance) Len() int           { return len(a) }
func (a ByDistance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDistance) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

func main() {
    var queryX, queryY float64
    var k int

    fmt.Print("Enter the X and Y coordinates of the query point (separated by space): ")
    _, err := fmt.Scanf("%f %f", &queryX, &queryY)
    if err != nil {
        log.Fatalf("Failed to read coordinates: %v", err)
    }

    fmt.Print("Enter the value of k (number of nearest neighbors to find): ")
    _, err = fmt.Scanln(&k)
    if err != nil {
        log.Fatalf("Failed to read k value: %v", err)
    }

    queryPoint := &knn.DataPoint{Coordinates: []float32{float32(queryX), float32(queryY)}}
    var allNeighbors []Neighbor
    var mu sync.Mutex
    var wg sync.WaitGroup

    for i := 0; i < numServers; i++ {
        port := basePort + i
        wg.Add(1)

        go func(port int) {
            defer wg.Done()

            conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
            if err != nil {
                log.Printf("did not connect to server %d: %v", port, err)
                return
            }
            defer conn.Close()

            client := knn.NewKNNServiceClient(conn)

            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            defer cancel()

            stream, err := client.FindKNearestNeighbors(ctx, &knn.KNNRequest{QueryPoint: queryPoint, K: int32(k)})
            if err != nil {
                log.Printf("could not call FindKNearestNeighbors on server %d: %v", port, err)
                return
            }

            for {
                neighbor, err := stream.Recv()
                if err == io.EOF {
                    break
                }
                if err != nil {
                    log.Printf("error receiving neighbors from server %d: %v", port, err)
                    return
                }

                mu.Lock()
                allNeighbors = append(allNeighbors, Neighbor{
                    Point:    neighbor.Point,
                    Distance: neighbor.Distance,
                })
                mu.Unlock()
            }
        }(port)
    }

    wg.Wait()

    sort.Sort(ByDistance(allNeighbors))

    if len(allNeighbors) < k {
        log.Fatalf("Error: Queried k (%d) is greater than the number of available points (%d).", k, len(allNeighbors))
    }

    allNeighbors = allNeighbors[:k]

    fmt.Println("Global K Nearest Neighbors:")
    for _, neighbor := range allNeighbors {
        fmt.Printf("Point: %v, Distance: %f\n", neighbor.Point.Coordinates, neighbor.Distance)
    }
}
