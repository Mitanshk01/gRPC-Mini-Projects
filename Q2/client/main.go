package main

import (
    "container/heap"
    "context"
    "fmt"
    "io"
    "log"
    "sync"
    "time"

    "google.golang.org/grpc"
    pb "github.com/Mitanshk01/DS_HW4/Q2/protofiles"
)

const (
    numServers = 5
    basePort   = 50051
)

type Neighbor struct {
    Point    *pb.DataPoint
    Distance float32
}

// MaxHeap interface implementation
type MaxHeap []Neighbor

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].Distance > h[j].Distance }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x interface{}) {
    *h = append(*h, x.(Neighbor))
}

func (h *MaxHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

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

    queryPoint := &pb.DataPoint{Coordinates: []float32{float32(queryX), float32(queryY)}}
    neighborHeap := &MaxHeap{}
    heap.Init(neighborHeap)

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

            client := pb.NewKNNServiceClient(conn)

            ctx, cancel := context.WithTimeout(context.Background(), time.Second)
            defer cancel()

            stream, err := client.FindKNearestNeighbors(ctx, &pb.KNNRequest{QueryPoint: queryPoint, K: int32(k)})
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
                if neighborHeap.Len() < k {
                    heap.Push(neighborHeap, Neighbor{
                        Point:    neighbor.Point,
                        Distance: neighbor.Distance,
                    })
                } else if neighbor.Distance < (*neighborHeap)[0].Distance {
                    heap.Pop(neighborHeap)
                    heap.Push(neighborHeap, Neighbor{
                        Point:    neighbor.Point,
                        Distance: neighbor.Distance,
                    })
                }
                mu.Unlock()
            }
        }(port)
    }

    wg.Wait()

    if neighborHeap.Len() < k {
        log.Fatalf("Error: Queried k (%d) is greater than the number of available points (%d).", k, neighborHeap.Len())
    }

    fmt.Println("Global K Nearest Neighbors:")
    results := make([]Neighbor, neighborHeap.Len())
    for i := len(results) - 1; i >= 0; i-- {
        results[i] = heap.Pop(neighborHeap).(Neighbor)
    }
    for _, neighbor := range results {
        fmt.Printf("Point: %v, Distance: %f\n", neighbor.Point.Coordinates, neighbor.Distance)
    }
}