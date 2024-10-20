# Homework-4: gRPC

```
Names : Mitansh Kayathwal, Pradeep Mishra
Roll Nos: 2021101026, 2023801013
Branch : CSE, PhD
Course : Distributed Systems, Monsoon '24
```

# **_Directory Structure_**

```
📁 Q4
├── 📁 client
│   └── 📄 main.go
├── 📁 logs
│   └── 📄 logger.go
│   └── 📄 logs.txt
├── 📁 protofiles
│   └── 📄 document_grpc.pb.go
│   └── 📄 document.pb.go
│   └── 📄 document.proto
├── 📁 server
│   └── 📄 main.go
├── 📄 go.mod
├── 📄 go.sum
├── 📄 Makefile
├── 📄 README.md
├── 📄 Report.pdf
├── 📄 toBuild.sh
```

# Implementation Quirks

## Overview

The Live Document system is implemented using gRPC for efficient client-server communication. The system allows multiple clients to collaboratively edit a shared document in real-time. Here's a brief overview of the implementation:

1. **Protocol Buffers**: The `document.proto` file defines the service and message types for the Live Document system.

2. **Server**: The server maintains the current state of the document and handles client requests for updates and modifications.

3. **Client**: Clients can connect to the server, request the current document state, and send updates to modify the document.

4. **Bi-directional Streaming**: The system uses bi-directional gRPC streaming to enable real-time updates between the server and connected clients.

5. **Logging**: A logging system is implemented to track all operations and changes made to the document.

## Components

1. **Server (server/main.go)**:
   - Implements the gRPC service defined in the proto file.
   - Manages the document state and client connections.
   - Broadcasts updates to all connected clients.

2. **Client (client/main.go)**:
   - Connects to the server using gRPC.
   - Sends document modification requests to the server.
   - Receives and displays real-time updates from the server.

3. **Proto Files (protofiles/document.proto)**:
   - Defines the service interface and message types for the Live Document system.

4. **Logger (logger/logger.go)**:
   - Implements logging functionality to record all document operations.


## Running the Code

To run the Live Document editor, follow these steps (from the root directory Q4/):

1. Make the build script executable and run it:
   ```
   $ chmod +x toBuild.sh
   $ ./toBuild.sh
   ```

   This script will clean up previous builds and generate the necessary proto files.


2. Start the server:
   ```
   $ make server
   ```

3. Start the logger:
   ```
   $ make logger
   ```

4. In a separate terminal, start the 2 clients (can be extended to more clients, this is just for testing purposes, I've tested on more than 2 clients and the system works smoothly in that scenario as well):
   ```
   $ make client1
   ```

   ```
   $ make client2
   ```

5. Follow the prompts in the client to interact with the Live Document system.

## Makefile Usage

The Makefile provides several commands to simplify development and testing.

Note that to run the program, you also need to run the following commands in order given below (from the root directory Q4/):

- `make clean`: Removes generated files (proto files).
- `make proto`: Generates Go code from the proto files.
- `make server`: Builds and runs the server.
- `make logger`: Builds and runs the logger.
- `make client1`: Builds and runs the client1 (port: 8080).
- `make client2`: Builds and runs the client1 (port: 8081).

To use these commands, simply run `make <command>` in the terminal (from the root directory Q4/).