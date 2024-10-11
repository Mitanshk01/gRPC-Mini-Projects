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
DRIVER1_PORT = 50054
DRIVER2_PORT = 50055
RIDER1_PORT = 50056
RIDER2_PORT = 50057

.PHONY: etcd proto server1 server2 server3 driver1 driver2 rider1 rider2 clean

etcd:
	go run $(SERVER_DIR)/etcd

proto:
	protoc $(GO_FLAGS) $(PROTO_FILE)

server1:
	go run $(SERVER_DIR)/main.go --port=$(SERVER1_PORT)

server2:
	go run $(SERVER_DIR)/main.go --port=$(SERVER2_PORT)

server3:
	go run $(SERVER_DIR)/main.go --port=$(SERVER3_PORT)

driver1:
	go run $(CLIENT_DIR)/driver.go --port=$(DRIVER1_PORT)

driver2:
	go run $(CLIENT_DIR)/driver.go --port=$(DRIVER2_PORT)

rider1:
	go run $(CLIENT_DIR)/rider.go --port=$(RIDER1_PORT)

rider2:
	go run $(CLIENT_DIR)/rider.go --port=$(RIDER2_PORT)

clean:
	rm -f $(PROTO_DIR)/*.pb.go
