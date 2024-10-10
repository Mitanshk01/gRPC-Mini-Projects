package main

import (
    "context"
    "encoding/json" // Import to handle JSON parsing
    "flag"
    "fmt"
    "log"
    "os"
    "time"

    "google.golang.org/grpc"
    pb "github.com/Mitanshk01/DS_HW4/Q4/protofiles"
)

// LoggerClient struct
type LoggerClient struct {
    logFile *os.File
}

// NewLoggerClient initializes the logger
func NewLoggerClient() *LoggerClient {
    logFile, err := os.OpenFile("./logs/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatalf("Error opening log file: %v", err)
    }

    return &LoggerClient{logFile: logFile}
}

// StreamDocumentLogs receives log data from the server
func (lc *LoggerClient) StreamDocumentLogs(client pb.CollaborativeDocumentServiceClient) error {
    stream, err := client.StreamDocumentLogs(context.Background(), &pb.EmptyMessage{})
    if err != nil {
        return fmt.Errorf("failed to establish log stream: %v", err)
    }

    for {
        change, err := stream.Recv()
        if err != nil {
            return fmt.Errorf("error receiving log stream: %v", err)
        }

        // Parse and format the timestamp
        ts, err := time.Parse(time.RFC3339, change.Timestamp)
        if err != nil {
            log.Printf("Error parsing timestamp: %v, using raw value.", err)
            ts = time.Now()
        }
        formattedTimestamp := ts.Format("2006-01-02 15:04:05")

        var logLine string

        switch change.ChangeType {
        case "edit":
            var contentData struct {
                DeleteContent string `json:"deleteContent"`
                AddContent    string `json:"addContent"`
            }
            if err := json.Unmarshal([]byte(change.Content), &contentData); err == nil {
                logLine = fmt.Sprintf("[%s] ClientID: %s | ChangeType: %s | Position: %d | OldContent: %s | NewContent: %s\n",
                    formattedTimestamp, change.ClientId, change.ChangeType, change.Position, contentData.DeleteContent, contentData.AddContent)
            } else {
                log.Printf("Error parsing JSON content for edit: %v", err)
                logLine = fmt.Sprintf("[%s] ClientID: %s | ChangeType: %s | Position: %d | Content: %s\n",
                    formattedTimestamp, change.ClientId, change.ChangeType, change.Position, change.Content)
            }
        case "add", "delete":
            logLine = fmt.Sprintf("[%s] ClientID: %s | ChangeType: %s | Position: %d | Content: %s\n",
                formattedTimestamp, change.ClientId, change.ChangeType, change.Position, change.Content)
        default:
            logLine = fmt.Sprintf("[%s] ClientID: %s | ChangeType: %s | Position: %d | Content: %s\n",
                formattedTimestamp, change.ClientId, change.ChangeType, change.Position, change.Content)
        }

        if _, err := lc.logFile.WriteString(logLine); err != nil {
            return fmt.Errorf("failed to write log: %v", err)
        }
    }
}

func main() {
    port := flag.String("port", "50051", "Port to connect to the server")
    flag.Parse()

    conn, err := grpc.Dial("localhost:"+*port, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect to server: %v", err)
    }
    defer conn.Close()

    client := pb.NewCollaborativeDocumentServiceClient(conn)
    loggerClient := NewLoggerClient()

    if err := loggerClient.StreamDocumentLogs(client); err != nil {
        log.Fatalf("Error streaming logs: %v", err)
    }

    defer loggerClient.logFile.Close()
}
