package main

import (
    "bufio"
    "context"
    "flag"
    "fmt"
    "log"
    "os"
    "strings"
    "time"
    "strconv"

    "google.golang.org/grpc"
    pb "github.com/Mitanshk01/DS_HW4/Q1/protofiles"
)

func displayCommands() {
    fmt.Println("Available commands:")
    fmt.Println("1. labyrinth_info - Get labyrinth dimensions")
    fmt.Println("2. player_status - Get player score, health, and position")
    fmt.Println("3. move [DIRECTION] - Move the player (DIRECTION: U, D, L, R)")
    fmt.Println("4. revelio [X] [Y] [TileType] - Use Revelio spell at given position and tile type")
    fmt.Println("5. bombarda [X] [Y] - Use Bombarda spell at given position")
    fmt.Println("6. exit - Exit the game")
}

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

func useRevelio(client pb.LabyrinthServiceClient, x, y uint32, tileType string) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    revelioRequest := &pb.RevelioRequest{
        TargetPosition: &pb.Position{PositionX: x, PositionY: y},
        TileType:       tileType,
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

func useBombarda(client pb.LabyrinthServiceClient, x, y uint32) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    bombardaStream, err := client.Bombarda(ctx)
    if err != nil {
        log.Fatalf("Could not start Bombarda: %v", err)
    }

    bombardaRequest := &pb.BombardaRequest{
        TargetPosition: &pb.Position{PositionX: x, PositionY: y},
    }

    if err := bombardaStream.Send(bombardaRequest); err != nil {
        log.Fatalf("Bombarda send error: %v", err)
    }

    fmt.Printf("Bombarda spell cast at position (%d, %d)!\n", x, y)
    bombardaStream.CloseSend()
}

func main() {
    serverAddress := flag.String("server", "localhost:50051", "The server address in the format of host:port")
    flag.Parse()

    conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    labyrinthClient := pb.NewLabyrinthServiceClient(conn)
    playerClient := pb.NewPlayerServiceClient(conn)

    reader := bufio.NewReader(os.Stdin)
    displayCommands()

    for {
        fmt.Print("\nEnter command: ")
        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(input)
        parts := strings.Fields(input)

        if len(parts) == 0 {
            continue
        }

        command := parts[0]

        switch command {
        case "labyrinth_info":
            getLabyrinthInfo(labyrinthClient)

        case "player_status":
            getPlayerStatus(playerClient)

        case "move":
            if len(parts) < 2 {
                fmt.Println("Usage: move [DIRECTION]")
                continue
            }
            direction := strings.ToUpper(parts[1])
            if direction != "U" && direction != "D" && direction != "L" && direction != "R" {
                fmt.Println("Invalid direction! Use U (up), D (down), L (left), or R (right).")
                continue
            }
            registerMove(playerClient, direction)

        case "revelio":
            if len(parts) < 4 {
                fmt.Println("Usage: revelio [X] [Y] [TileType]")
                continue
            }
            x := parseToUInt32(parts[1])
            y := parseToUInt32(parts[2])
            tileType := parts[3]
            useRevelio(labyrinthClient, x, y, tileType)

        case "bombarda":
            if len(parts) < 3 {
                fmt.Println("Usage: bombarda [X] [Y]")
                continue
            }
            x := parseToUInt32(parts[1])
            y := parseToUInt32(parts[2])
            useBombarda(labyrinthClient, x, y)

        case "exit":
            fmt.Println("Exiting the game...")
            return

        default:
            fmt.Println("Unknown command! Please try again.")
            displayCommands()
        }
    }
}

func parseToUInt32(s string) uint32 {
    val, err := strconv.ParseUint(s, 10, 32)
    if err != nil {
        fmt.Printf("Error parsing '%s' to uint32: %v\n", s, err)
        return 0
    }
    return uint32(val)
}