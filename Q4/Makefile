PROTO_DIR = protofiles
CLIENT_DIR = client
SERVER_DIR = server
LOGS_DIR = logs

PROTO_FILE = $(PROTO_DIR)/document.proto
PROTO_OUT_DIR = .

GO_FLAGS = --go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
           --go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative

SERVER_PORT = 50051
CLIENT1_PORT = 8080
CLIENT2_PORT = 8081

.PHONY: all proto server client1 client2 clean

all: proto server client1 client2

proto:
	protoc $(GO_FLAGS) $(PROTO_FILE)

server:
	go run $(SERVER_DIR)/main.go --port=$(SERVER_PORT)

logger:
	go run $(LOGS_DIR)/logger.go --port=$(SERVER_PORT)

client1:
	go run $(CLIENT_DIR)/main.go --port=$(CLIENT1_PORT)

client2:
	go run $(CLIENT_DIR)/main.go --port=$(CLIENT2_PORT)

clean:
	rm -f $(PROTO_DIR)/*.pb.go
