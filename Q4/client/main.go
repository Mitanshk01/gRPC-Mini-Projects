package main

import (
	"context"
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	pb "github.com/Mitanshk01/DS_HW4/Q4/protofiles"
	"github.com/google/uuid"
)

var (
	grpcAddr   = "localhost:50051"
	clientID   = uuid.New().String()
	docContent = ""
	mu         sync.Mutex
	clients    = make(map[*websocket.Conn]bool)
	grpcStream pb.CollaborativeDocumentService_SyncDocumentChangesClient
)

type changeData struct {
	DeleteContent  string `json:"deleteContent"`
	DeletePosition int    `json:"deletePosition"`
	AddContent     string `json:"addContent"`
	AddPosition    int    `json:"addPosition"`
	ChangeType     string `json:"change_type"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	for {
		var msg changeData
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading JSON: %v", err)
			delete(clients, conn)
			break
		}

		handleDocumentEdit(msg)
	}
}

func handleDocumentEdit(msg changeData) {
	mu.Lock()
	defer mu.Unlock()

	var change *pb.DocumentChange

	switch msg.ChangeType {
	case "add":
		change = &pb.DocumentChange{
			ClientId:   clientID,
			Content:    msg.AddContent,
			Position:   int32(msg.AddPosition),
			ChangeType: "add",
		}
		docContent = docContent[:msg.AddPosition] + msg.AddContent + docContent[msg.AddPosition:]

	case "delete":
		change = &pb.DocumentChange{
			ClientId:   clientID,
			Content:    msg.DeleteContent,
			Position:   int32(msg.DeletePosition),
			ChangeType: "delete",
		}
		docContent = docContent[:msg.DeletePosition - 1] + docContent[msg.DeletePosition+len(msg.DeleteContent)-1:]

	case "replace":
		jsonContent, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshaling replace change: %v", err)
			return
		}
		change = &pb.DocumentChange{
			ClientId:   clientID,
			Content:    string(jsonContent),
			Position:   int32(msg.DeletePosition),
			ChangeType: "edit",
		}
		docContent = docContent[:msg.DeletePosition] + msg.AddContent + docContent[msg.DeletePosition+len(msg.DeleteContent):]
	}

	change.Timestamp = time.Now().Format(time.RFC3339)

	if err := grpcStream.Send(change); err != nil {
		log.Printf("Error sending change to gRPC server: %v", err)
	}
}


func SyncDocumentChanges(conn pb.CollaborativeDocumentServiceClient) {
	var err error
	grpcStream, err = conn.SyncDocumentChanges(context.Background())
	if err != nil {
		log.Fatalf("Error connecting to stream: %v", err)
	}

	initialChange := &pb.DocumentChange{
		ClientId:   clientID,
		Content:    "",
		ChangeType: "initial",
		Position:   0,
	}

	if err := grpcStream.Send(initialChange); err != nil {
		log.Fatalf("Error sending initial request: %v", err)
	}

	go func() {
		for {
			in, err := grpcStream.Recv()
			if err != nil {
				log.Fatalf("Error receiving changes: %v", err)
			}
			
			mu.Lock()
			switch in.ChangeType {
			case "initial":
				docContent = in.Content
			case "add":
				docContent = docContent[:in.Position] + in.Content + docContent[in.Position:]
			case "delete":
				log.Printf("%d, %d", docContent, in.Position, in.Content)
				docContent = docContent[:in.Position - 1] + docContent[in.Position+int32(len(in.Content))-1:]
			case "edit":
				var change changeData
				if err := json.Unmarshal([]byte(in.Content), &change); err != nil {
					continue
				}
				docContent = docContent[:change.DeletePosition] + change.AddContent + docContent[change.DeletePosition+len(change.DeleteContent):]
			}
			mu.Unlock()

			broadcastChange(in)
		}
	}()
}

func broadcastChange(change *pb.DocumentChange) {
	for client := range clients {
		err := client.WriteJSON(map[string]string{
			"Content": docContent,
		})
		if err != nil {
			log.Printf("Error sending change to client: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}


func renderHTML(w http.ResponseWriter, r *http.Request) {
	const htmlContent = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Collaborative Document Editor</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				margin: 0;
				padding: 20px;
				background-color: #f4f4f4;
			}
			#editor {
				width: 100%;
				height: 400px;
				border: 1px solid #ccc;
				padding: 10px;
				font-size: 16px;
				background-color: #fff;
			}
		</style>
	</head>
	<body>
		<h1>Collaborative Document Editor</h1>
		<textarea id="editor">{{.}}</textarea>
		<script>
		const editor = document.getElementById('editor');
		let lastContent = editor.value;
		let isRemoteUpdate = false;
		let lastSelectionStart = editor.selectionStart;
		let lastSelectionEnd = editor.selectionEnd;

		const socket = new WebSocket('ws://' + window.location.host + '/ws');

		socket.onopen = () => {
			console.log('WebSocket connection established');
		};

		socket.onmessage = (event) => {
			const data = JSON.parse(event.data);
			isRemoteUpdate = true;
		
			const currentSelectionStart = editor.selectionStart;
			const currentSelectionEnd = editor.selectionEnd;
			const previousContent = editor.value;
			const previousContentLength = previousContent.length;
		
			const newContent = data.Content;
			const newContentLength = newContent.length;
		
			let diffPosition = 0;
			while (
				diffPosition < previousContentLength &&
				diffPosition < newContentLength &&
				previousContent[diffPosition] === newContent[diffPosition]
			) {
				diffPosition++;
			}
		
			editor.value = newContent;
		
			if (diffPosition <= currentSelectionStart) {
				const lengthDifference = newContentLength - previousContentLength;
				editor.selectionStart = currentSelectionStart + lengthDifference;
				editor.selectionEnd = currentSelectionEnd + lengthDifference;
			} else {
				editor.selectionStart = currentSelectionStart;
				editor.selectionEnd = currentSelectionEnd;
			}
		
			lastContent = editor.value;
			isRemoteUpdate = false;
		};

		editor.addEventListener('beforeinput', (e) => {
			lastSelectionStart = editor.selectionStart;
			lastSelectionEnd = editor.selectionEnd;
		});

		editor.addEventListener('input', () => {
			if (isRemoteUpdate) return;

			const currentContent = editor.value;
			const selectionStart = editor.selectionStart;
			const selectionEnd = editor.selectionEnd;

			if (lastSelectionStart !== lastSelectionEnd) {
				const deletedContent = lastContent.slice(lastSelectionStart, lastSelectionEnd);
				const addedContent = currentContent.slice(lastSelectionStart, selectionStart);
				
				console.log("Sending edit request");
				socket.send(JSON.stringify({
					deleteContent: deletedContent,
					deletePosition: lastSelectionStart,
					addContent: addedContent,
					addPosition: lastSelectionStart,
					change_type: 'replace'
				}));
			} else if (currentContent.length > lastContent.length) {
				const addedContent = currentContent.slice(lastSelectionStart, selectionStart);
				
				console.log("Sending add request");
				socket.send(JSON.stringify({
					deleteContent: '',
					deletePosition: lastSelectionStart,
					addContent: addedContent,
					addPosition: lastSelectionStart,
					change_type: 'add'
				}));
			} else if (currentContent.length < lastContent.length) {
				const deletedContent = lastContent.slice(selectionStart, selectionStart + (lastContent.length - currentContent.length));
				console.log("Sending delete request", lastSelectionStart);
				socket.send(JSON.stringify({
					deleteContent: deletedContent,
					deletePosition: lastSelectionStart,
					addContent: '',
					addPosition: lastSelectionStart,
					change_type: 'delete'
				}));
			}

			lastContent = currentContent;
		});
		</script>
	</body>
	</html>
	`

	tmpl, err := template.New("editor").Parse(htmlContent)
	if err != nil {
		http.Error(w, "Could not parse template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, docContent)
	if err != nil {
		http.Error(w, "Could not render template", http.StatusInternalServerError)
	}
}

func main() {
	port := flag.String("port", "", "Port to connect to the server (e.g., 50051)")
	flag.Parse()

	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()
	log.Println("Client ID:", clientID)

	client := pb.NewCollaborativeDocumentServiceClient(conn)

	go SyncDocumentChanges(client)

	http.HandleFunc("/", renderHTML)
	http.HandleFunc("/ws", handleWebSocket)

	go func() {
		log.Printf("Server running on port %s", *port)
		if err := http.ListenAndServe("127.0.0.1:"+*port, nil); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down server...")
}