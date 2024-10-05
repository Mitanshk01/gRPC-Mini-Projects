package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "google.golang.org/grpc"
    pb "github.com/Mitanshk01/DS_HW4/Q1/protofiles"
)

func getLabyrinthInfo(client pb.LabyrinthServiceClient) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    response, err := client.GetLabyrinthInfo(ctx, &pb.EmptyMessage{})
    if err != nil {
        log.Fatalf("Could not get labyrinth info: %v", err)
    }
    fmt.Printf("Labyrinth dimensions - Width: %d, Height: %d\n", response.Width, response.Height)
}

func getPlayerStatus(client pb.PlayerServiceClient) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    response, err := client.GetPlayerStatus(ctx, &pb.EmptyMessage{})
    if err != nil {
        log.Fatalf("Could not get player status: %v", err)
    }
    fmt.Printf("Player Status - Score: %d, Health: %d, Position: (%d, %d)\n",
        response.Score, response.HealthPoints, response.Position.PositionX, response.Position.PositionY)
}

func registerMove(client pb.PlayerServiceClient, direction string) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    moveRequest := &pb.MoveRequest{Direction: direction}
    response, err := client.RegisterMove(ctx, moveRequest)
    if err != nil {
        log.Fatalf("Could not register move: %v", err)
    }

    switch response.Result {
    case pb.MoveResult_SUCCESS:
        fmt.Println("Move successful!")
    case pb.MoveResult_FAILURE:
        fmt.Println("Move failed!")
    case pb.MoveResult_PLAYER_DEAD:
        fmt.Println("Player is dead!")
    case pb.MoveResult_VICTORY:
        fmt.Println("Victory!")
    }
}

func useRevelio(client pb.LabyrinthServiceClient) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    revelioRequest := &pb.RevelioRequest{
        TargetPosition: &pb.Position{PositionX: 1, PositionY: 1},
        TileType:       "C",
    }

    stream, err := client.Revelio(ctx, revelioRequest)
    if err != nil {
        log.Fatalf("Revelio failed: %v", err)
    }

    fmt.Println("Revelio results:")
    for {
        position, err := stream.Recv()
        if err != nil {
            break
        }
        fmt.Printf("Tile found at position (%d, %d)\n", position.PositionX, position.PositionY)
    }
}

func useBombarda(client pb.LabyrinthServiceClient) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    bombardaStream, err := client.Bombarda(ctx)
    if err != nil {
        log.Fatalf("Could not start Bombarda: %v", err)
    }

    bombardaRequest := &pb.BombardaRequest{
        TargetPosition: &pb.Position{PositionX: 2, PositionY: 2}, // Example position
    }

    if err := bombardaStream.Send(bombardaRequest); err != nil {
        log.Fatalf("Bombarda send error: %v", err)
    }

    fmt.Println("Bombarda spell cast at position (2, 2)!")
    bombardaStream.CloseSend()
}

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    labyrinthClient := pb.NewLabyrinthServiceClient(conn)
    playerClient := pb.NewPlayerServiceClient(conn)

    getLabyrinthInfo(labyrinthClient)

    getPlayerStatus(playerClient)

    registerMove(playerClient, "R")

    useRevelio(labyrinthClient)

    useBombarda(labyrinthClient)
}
