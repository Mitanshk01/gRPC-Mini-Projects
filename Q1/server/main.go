package main

import (
    "bufio"
    "context"
    "errors"
    "flag"
    "log"
    "net"
    "os"
    "strings"

    "google.golang.org/grpc"
	pb "github.com/Mitanshk01/DS_HW4/Q1/protofiles"
)

type LabyrinthServer struct {
    pb.UnimplementedLabyrinthServiceServer
    labyrinth [][]Tile
    player    Player
}

type PlayerServer struct {
    pb.UnimplementedPlayerServiceServer
    labyrinthServer *LabyrinthServer
}

type Tile struct {
    Type string  // "EMPTY", "COIN", "WALL"
}

type Player struct {
    PositionX         uint32
    PositionY         uint32
    Health            uint32
    Score             uint32
    Spells            uint32
}

func isPortAvailable(port string) bool {
    ln, err := net.Listen("tcp", ":"+port)
    if err != nil {
        return false
    }
    ln.Close()
    return true
}

func initializeLabyrinthFromFile(filePath string) [][]Tile {
    file, err := os.Open(filePath)
    if err != nil {
        log.Fatalf("Failed to open input file: %v", err)
    }
    defer file.Close()

    var labyrinth [][]Tile
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := scanner.Text()
        row := []Tile{}

        tileSymbols := strings.Split(line, " ")
        for _, symbol := range tileSymbols {
            var tileType string
            switch symbol {
            case "E":
                tileType = "EMPTY"
            case "C":
                tileType = "COIN"
            case "W":
                tileType = "WALL"
            default:
                log.Fatalf("Invalid tile symbol in labyrinth file: %v", symbol)
            }
            row = append(row, Tile{Type: tileType})
        }
        labyrinth = append(labyrinth, row)
    }

    if err := scanner.Err(); err != nil {
        log.Fatalf("Error reading labyrinth file: %v", err)
    }

    return labyrinth
}

func (s *LabyrinthServer) GetLabyrinthInfo(ctx context.Context, req *pb.EmptyMessage) (*pb.LabyrinthInfoResponse, error) {
    return &pb.LabyrinthInfoResponse{
        Width:  uint32(len(s.labyrinth[0])),
        Height: uint32(len(s.labyrinth)),
    }, nil
}

func (s *LabyrinthServer) Revelio(req *pb.RevelioRequest, stream pb.LabyrinthService_RevelioServer) error {
    if s.player.Spells <= 0 {
        return errors.New("No spells remaining")
    }

    // log.Printf("Current spells: %d", s.player.Spells)

    targetX, targetY := req.TargetPosition.PositionX, req.TargetPosition.PositionY
    
    if targetX >= uint32(len(s.labyrinth[0])) || targetY >= uint32(len(s.labyrinth)) {
        return errors.New("Invalid target position")
    }

    startX := int(targetX) - 1
    startY := int(targetY) - 1
    endX := int(targetX) + 1
    endY := int(targetY) + 1

    if startX < 0 {
        startX = 0
    }
    if startY < 0 {
        startY = 0
    }
    if endX >= len(s.labyrinth[0]) {
        endX = len(s.labyrinth[0]) - 1
    }
    if endY >= len(s.labyrinth) {
        endY = len(s.labyrinth) - 1
    }

    for y := startY; y <= endY; y++ {
        for x := startX; x <= endX; x++ {
            if x >= 0 && x < len(s.labyrinth[0]) && y >= 0 && y < len(s.labyrinth) {
                tile := s.labyrinth[y][x]
                if tile.Type == req.TileType || req.TileType == "" {
                    if err := stream.Send(&pb.Position{
                        PositionX: uint32(x),
                        PositionY: uint32(y),
                    }); err != nil {
                        return err
                    }
                }
            }
        }
    }

    s.player.Spells--
    return nil
}

func (s *LabyrinthServer) Bombarda(stream pb.LabyrinthService_BombardaServer) error {
    if s.player.Spells <= 0 {
        return errors.New("No spells remaining")
    }

    var requests []*pb.BombardaRequest
    for {
        req, err := stream.Recv()
        if err != nil {
            if len(requests) != 3 {
                return errors.New("Invalid number of points received, expected exactly 3")
            }
            break
        }

        if req.TargetPosition.PositionX >= uint32(len(s.labyrinth[0])) || 
           req.TargetPosition.PositionY >= uint32(len(s.labyrinth)) {
            return errors.New("Invalid target position")
        }

        requests = append(requests, req)
    }

    for _, req := range requests {
        x, y := req.TargetPosition.PositionX, req.TargetPosition.PositionY

        startX := int(x) - 1
        startY := int(y) - 1
        endX := int(x) + 1
        endY := int(y) + 1

        if startX < 0 {
            startX = 0
        }
        if startY < 0 {
            startY = 0
        }
        if endX >= len(s.labyrinth[0]) {
            endX = len(s.labyrinth[0]) - 1
        }
        if endY >= len(s.labyrinth) {
            endY = len(s.labyrinth) - 1
        }

        for j := startY; j <= endY; j++ {
            for i := startX; i <= endX; i++ {
                s.labyrinth[j][i] = Tile{Type: "EMPTY"}
            }
        }
    }

    s.player.Spells--
    return stream.SendAndClose(&pb.EmptyMessage{})
}

func (s *PlayerServer) GetPlayerStatus(ctx context.Context, req *pb.EmptyMessage) (*pb.PlayerStatusResponse, error) {
    return &pb.PlayerStatusResponse{
        Score:        s.labyrinthServer.player.Score,
        HealthPoints: s.labyrinthServer.player.Health,
        Position: &pb.Position{
            PositionX: s.labyrinthServer.player.PositionX,
            PositionY: s.labyrinthServer.player.PositionY,
        },
    }, nil
}

func handleWallCollision(s *PlayerServer) (*pb.MoveResponse, error) {
    s.labyrinthServer.player.Health--
    if s.labyrinthServer.player.Health == 0 {
        return &pb.MoveResponse{Result: pb.MoveResult_PLAYER_DEAD}, nil
    }
    return &pb.MoveResponse{Result: pb.MoveResult_FAILURE}, nil
}

func (s *PlayerServer) RegisterMove(ctx context.Context, req *pb.MoveRequest) (*pb.MoveResponse, error) {
    direction := req.Direction
    x, y := s.labyrinthServer.player.PositionX, s.labyrinthServer.player.PositionY

    if x == uint32(len(s.labyrinthServer.labyrinth[0]))-1 && y == uint32(len(s.labyrinthServer.labyrinth))-1 {
        return &pb.MoveResponse{Result: pb.MoveResult_VICTORY}, nil
    }

    if s.labyrinthServer.player.Health == 0 {
        return &pb.MoveResponse{Result: pb.MoveResult_PLAYER_DEAD}, nil
    }

    switch direction {
    case "U":
        if y > 0 {
            y--
        } else {
            return handleWallCollision(s)
        }
    case "D":
        if y < uint32(len(s.labyrinthServer.labyrinth))-1 {
            y++
        } else {
            return handleWallCollision(s)
        }
    case "L":
        if x > 0 {
            x--
        } else {
            return handleWallCollision(s)
        }
    case "R":
        if x < uint32(len(s.labyrinthServer.labyrinth[0]))-1 {
            x++
        } else {
            return handleWallCollision(s)
        }
    default:
        return &pb.MoveResponse{Result: pb.MoveResult_FAILURE}, nil
    }

    tile := s.labyrinthServer.labyrinth[y][x]
    if tile.Type == "WALL" {
        return handleWallCollision(s)
    } else if tile.Type == "COIN" {
        s.labyrinthServer.player.Score++
        s.labyrinthServer.labyrinth[y][x] = Tile{Type: "EMPTY"}
    }

    s.labyrinthServer.player.PositionX = x
    s.labyrinthServer.player.PositionY = y

    if x == uint32(len(s.labyrinthServer.labyrinth[0]))-1 && y == uint32(len(s.labyrinthServer.labyrinth))-1 {
        return &pb.MoveResponse{Result: pb.MoveResult_VICTORY}, nil
    }

    return &pb.MoveResponse{Result: pb.MoveResult_SUCCESS}, nil
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

    listener, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    labyrinthFilePath := "./server/labyrinth.txt"
    labyrinth := initializeLabyrinthFromFile(labyrinthFilePath)

    labyrinthServer := &LabyrinthServer{
        labyrinth: labyrinth,
        player:    Player{PositionX: 0, PositionY: 0, Health: 3, Score: 0, Spells: 3},
    }

    grpcServer := grpc.NewServer()
    pb.RegisterLabyrinthServiceServer(grpcServer, labyrinthServer)
    pb.RegisterPlayerServiceServer(grpcServer, &PlayerServer{labyrinthServer: labyrinthServer})

    log.Println("Server is running on port 50051...")
    if err := grpcServer.Serve(listener); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
