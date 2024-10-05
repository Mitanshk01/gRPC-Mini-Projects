package main

import (
    "bufio"
    "context"
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
    targetX, targetY := req.TargetPosition.PositionX, req.TargetPosition.PositionY
    for y, row := range s.labyrinth {
        for x, tile := range row {
            if (tile.Type == req.TileType || req.TileType == "") && uint32(x) == targetX && uint32(y) == targetY {
                stream.Send(&pb.Position{
                    PositionX: uint32(x),
                    PositionY: uint32(y),
                })
            }
        }
    }
    return nil
}

func (s *LabyrinthServer) Bombarda(stream pb.LabyrinthService_BombardaServer) error {
    for {
        req, err := stream.Recv()
        if err != nil {
            return err
        }
        x, y := req.TargetPosition.PositionX, req.TargetPosition.PositionY
        if x < uint32(len(s.labyrinth[0])) && y < uint32(len(s.labyrinth)) {
            s.labyrinth[y][x] = Tile{Type: "EMPTY"}
        }
    }
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

    switch direction {
    case "U":
        if y > 0 {
            y--
        } else {
            handleWallCollision(s)
        }
    case "D":
        if y < uint32(len(s.labyrinthServer.labyrinth))-1 {
            y++
        } else {
            handleWallCollision(s)
        }
    case "L":
        if x > 0 {
            x--
        } else {
            handleWallCollision(s)
        }
    case "R":
        if x < uint32(len(s.labyrinthServer.labyrinth[0]))-1 {
            x++
        } else {
            handleWallCollision(s)
        }
    default:
        return &pb.MoveResponse{Result: pb.MoveResult_FAILURE}, nil
    }

    tile := s.labyrinthServer.labyrinth[y][x]
    if tile.Type == "WALL" {
        handleWallCollision(s)
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
    listener, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    labyrinthFilePath := "labyrinth.txt"
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
