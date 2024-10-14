PROTO_DIR = protofiles
CLIENT_DIR = client
SERVER_DIR = server
DATASET_DIR = dataset

PROTO_FILE = $(PROTO_DIR)/knn.proto
PROTO_OUT_DIR = .

GO_FLAGS = --go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
           --go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative

SERVER1_PORT = 50051
SERVER2_PORT = 50052
SERVER3_PORT = 50053
SERVER4_PORT = 50054
SERVER5_PORT = 50055

.PHONY: proto server client clean tests

proto:
	protoc $(GO_FLAGS) $(PROTO_FILE)

server1:
	go run $(SERVER_DIR)/main.go --port=$(SERVER1_PORT) --dataset=$(DATASET_DIR)/dataset_server_1.txt

server2:
	go run $(SERVER_DIR)/main.go --port=$(SERVER2_PORT) --dataset=$(DATASET_DIR)/dataset_server_2.txt

server3:
	go run $(SERVER_DIR)/main.go --port=$(SERVER3_PORT) --dataset=$(DATASET_DIR)/dataset_server_3.txt

server4:
	go run $(SERVER_DIR)/main.go --port=$(SERVER4_PORT) --dataset=$(DATASET_DIR)/dataset_server_4.txt

server5:
	go run $(SERVER_DIR)/main.go --port=$(SERVER5_PORT) --dataset=$(DATASET_DIR)/dataset_server_5.txt

client:
	go run $(CLIENT_DIR)/main.go

clean:
	rm -f $(PROTO_DIR)/*.pb.go
