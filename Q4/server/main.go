package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "log"
    "net"
    "sync"

    "google.golang.org/grpc"
    pb "github.com/Mitanshk01/DS_HW4/Q4/protofiles"
)

func isPortAvailable(port string) bool {
    ln, err := net.Listen("tcp", ":"+port)
    if err != nil {
        return false
    }
    ln.Close()
    return true
}

type changeData struct {
    DeleteContent  string `json:"deleteContent"`
    DeletePosition int    `json:"deletePosition"`
    AddContent     string `json:"addContent"`
    AddPosition    int    `json:"addPosition"`
    ChangeType     string `json:"change_type"`
}

type DocumentServer struct {
    pb.UnimplementedCollaborativeDocumentServiceServer
    documentContent string
    clients         map[string]pb.CollaborativeDocumentService_SyncDocumentChangesServer
    loggerStream    pb.CollaborativeDocumentService_StreamDocumentLogsServer
    mu              sync.Mutex
}

func NewDocumentServer() *DocumentServer {
    return &DocumentServer{
        documentContent: "",
        clients:         make(map[string]pb.CollaborativeDocumentService_SyncDocumentChangesServer),
    }
}

func (s *DocumentServer) SyncDocumentChanges(stream pb.CollaborativeDocumentService_SyncDocumentChangesServer) error {
    var clientID string

    for {
        change, err := stream.Recv()
        if err == io.EOF {
            log.Printf("Client %s disconnected.", clientID)
            break
        }
        if err != nil {
            return fmt.Errorf("error receiving stream from client: %v", err)
        }

        s.mu.Lock()

        if clientID == "" {
            clientID = change.ClientId
            s.clients[clientID] = stream
            log.Printf("Client %s connected.", clientID)

            initialChange := &pb.DocumentChange{
                ClientId:   clientID,
                Content:    s.documentContent,
                ChangeType: "initial",
                Position:   0,
                Timestamp:  change.Timestamp,
            }
            s.mu.Unlock()

            if err := stream.Send(initialChange); err != nil {
                log.Printf("Error sending initial content to client %s: %v", clientID, err)
                return err
            }

            s.mu.Lock()
        } else {
            switch change.ChangeType {
            case "add":
                pos := int(change.Position)
                s.documentContent = s.documentContent[:pos] + change.Content + s.documentContent[pos:]
            case "edit":
                var editChange changeData
                if err := json.Unmarshal([]byte(change.Content), &editChange); err != nil {
                    log.Printf("Error unmarshaling edit change: %v", err)
                    return err
                }

                deletePos := editChange.DeletePosition
                s.documentContent = s.documentContent[:deletePos] + s.documentContent[deletePos+len(editChange.DeleteContent):]

                addPos := editChange.AddPosition
                s.documentContent = s.documentContent[:addPos] + editChange.AddContent + s.documentContent[addPos:]

            case "delete":
                pos := int(change.Position) - 1
                s.documentContent = s.documentContent[:pos] + s.documentContent[pos+len(change.Content):]
            }

            for id, clientStream := range s.clients {
                if id != clientID {
                    go func(clientID string, clientStream pb.CollaborativeDocumentService_SyncDocumentChangesServer) {
                        if err := clientStream.Send(change); err != nil {
                            log.Printf("Error sending change to client %s: %v", clientID, err)
                        }
                    }(id, clientStream)
                }
            }

            if s.loggerStream != nil {
                go func() {
                    if err := s.loggerStream.Send(change); err != nil {
                        log.Printf("Error sending change to logger: %v", err)
                    }
                }()
            }
        }

        s.mu.Unlock()
    }

    s.mu.Lock()
    delete(s.clients, clientID)
    s.mu.Unlock()

    return nil
}

func (s *DocumentServer) StreamDocumentLogs(empty *pb.EmptyMessage, stream pb.CollaborativeDocumentService_StreamDocumentLogsServer) error {
    s.mu.Lock()
    s.loggerStream = stream
    s.mu.Unlock()

    log.Println("Logger client connected.")
    <-stream.Context().Done()
    log.Println("Logger client disconnected.")
    
    s.mu.Lock()
    s.loggerStream = nil
    s.mu.Unlock()

    return nil
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

    lis, err := net.Listen("tcp", ":"+*port)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    docServer := NewDocumentServer()
    pb.RegisterCollaborativeDocumentServiceServer(grpcServer, docServer)

    log.Printf("Starting collaborative document server on port %s...", *port)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
