PROTO_DIR = protofiles
CLIENT_DIR = client
SERVER_DIR = server

PROTO_FILE = $(PROTO_DIR)/ride.proto
PROTO_OUT_DIR = .

GO_FLAGS = --go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
           --go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative

SERVER1_PORT = 50051
SERVER2_PORT = 50052
SERVER3_PORT = 50053

.PHONY: proto server1 server2 server3 driver rider clean

# etcd_server:
# 	/$(SERVER_DIR)/etcd

proto:
	protoc $(GO_FLAGS) $(PROTO_FILE)

server1:
	go run $(SERVER_DIR)/main.go --port=$(SERVER1_PORT)

server2:
	go run $(SERVER_DIR)/main.go --port=$(SERVER2_PORT)

server3:
	go run $(SERVER_DIR)/main.go --port=$(SERVER3_PORT)

driver:
	go run $(CLIENT_DIR)/driver.go

rider:
	go run $(CLIENT_DIR)/rider.go

clean:
	rm -f $(PROTO_DIR)/*.pb.go
