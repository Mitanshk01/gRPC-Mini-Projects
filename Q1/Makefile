PROTO_DIR = protofiles
CLIENT_DIR = client
SERVER_DIR = server

PROTO_FILE_LABYRINTH = $(PROTO_DIR)/labyrinth.proto
PROTO_FILE_PLAYER = $(PROTO_DIR)/player.proto
PROTO_FILE_COMMON = $(PROTO_DIR)/common.proto
PROTO_OUT_DIR = .

GO_FLAGS = --go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
           --go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative

SERVER_PORT = 50051

.PHONY: proto server client clean tests

proto:
	protoc $(GO_FLAGS) $(PROTO_FILE_LABYRINTH) $(PROTO_FILE_PLAYER) $(PROTO_FILE_COMMON)

server:
	go run $(SERVER_DIR)/main.go --port=$(SERVER_PORT)

client:
	go run $(CLIENT_DIR)/main.go

tests:
	cd $(SERVER_DIR) && go test -v

clean:
	rm -f $(PROTO_DIR)/*.pb.go
