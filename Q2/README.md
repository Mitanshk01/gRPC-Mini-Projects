# Homework-4: gRPC

```
Names : Mitansh Kayathwal, Pradeep Mishra
Roll Nos: 2021101026, 2023801013
Branch : CSE, PhD
Course : Distributed Systems, Monsoon '24
```

# **_Directory Structure_**

```
📁 Q2
├── 📁 client
│   └── 📄 main.go
│   └── 📄 data_collection.go
├── 📁 dataset
│   └── 📄 dataset_server_1.txt
│   └── 📄 dataset_server_2.txt
│   └── 📄 dataset_server_3.txt
│   └── 📄 dataset_server_4.txt
│   └── 📄 dataset_server_5.txt 
│   └── 📄 dataset.txt 
│   └── 📄 generate_dataset.py 
│   └── 📄 partition_dataset.go
├── 📁 protofiles
│   └── 📄 knn_grpc.pb.go
│   └── 📄 knn.pb.go
│   └── 📄 knn.proto
├── 📁 server
│   └── 📄 main.go
├── 📄 go.mod
├── 📄 go.sum
├── 📄 Makefile
├── 📄 Report.pdf
├── 📄 README.md
```

# Implementation Details

## Overview
This part implements a distributed K-Nearest Neighbors (KNN) algorithm using gRPC for client-server communication. The system allows for parallel processing of KNN queries across multiple servers, each holding a portion of the dataset.

## Components

### Proto Files
- `knn.proto`: Defines the KNN service and related message types for distributed KNN implementation.

### Server
- `main.go`: Implements the KNN server, handling client queries and performing KNN calculations on its local dataset partition.

### Client
- `main.go`: Implements the client that sends KNN queries to multiple servers and aggregates results.
- `data_collection.go`: Implements the client with metrics logging to aggregate performance results.

### Dataset
- `dataset_server_*.txt`: Partitioned datasets for each server.
- `dataset.txt`: Complete dataset before partitioning.
- `generate_dataset.py`: Python script to generate the initial dataset.
- `partition_dataset.go`: Go script to partition the dataset among servers.

## Distributed KNN Approach

1. **Data Partitioning**: 
   - The complete dataset is partitioned among multiple servers using `partition_dataset.go`.
   - Each server loads its own partition of the dataset on startup.

2. **Client Query**:
   - The client sends a KNN query (point and K value) to all available servers.

3. **Server Processing**:
   - Each server performs KNN on its local dataset partition.
   - Servers return their top K nearest neighbors to the client.

4. **Result Aggregation**:
   - The client receives results from all servers.
   - It aggregates these results and performs a final KNN selection to get the global top K neighbors.

5. **Final Output**:
   - The client displays the final K nearest neighbors.

## Running the Code

To run the distributed KNN system, follow these steps (from the root directory Q2/):

1. Build the proto-files:
   ```
   make proto
   ```
   
   This script will clean up previous builds and generate the necessary proto files.

2. Start all the server instances (5 servers) in different terminals:
   ```
   $ make server1
   ```

   ```
   $ make server2
   ```

   ```
   $ make server3
   ```

   ```
   $ make server4
   ```

   ```
   $ make server5
   ```

3. In a separate terminal, start the client:
   ```
   $ make client
   ```

4. Follow the prompts in the client to input queries and receive KNN results.

## Makefile Usage

The Makefile provides several commands to simplify development and execution:

- `make clean`: Removes generated files (proto files and binaries).
- `make proto`: Generates Go code from the proto files.
- `make server1`, `make server2`, ...: Builds and runs individual server instances.
- `make client`: Builds and runs the client.
- `make generate_dataset`: Runs the Python script to generate the initial dataset.
- `make partition_dataset`: Runs the Go script to partition the dataset among servers.

To use these commands, run `make <command>` in the terminal from the root directory.

## Implementation Quirks and Optimizations

1. **Concurrent Server Queries**: The client sends queries to all servers concurrently using goroutines, improving overall query response time.

2. **Local KNN Optimization**: Each server performs KNN on its local dataset using an efficient algorithm (e.g., KD-tree or ball tree) for faster nearest neighbor search.

3. **Load Balancing**: The dataset partitioning ensures an even distribution of data points among servers, promoting balanced workload.

4. **Scalability**: The system is designed to easily accommodate additional servers by adding new partitions and updating the client's server list.

5. **Error Handling**: Robust error handling is implemented to manage scenarios such as server unavailability or network issues.

6. **Configurable K**: The system allows for dynamic K values, enabling flexibility in the number of nearest neighbors requested.

7. **Dataset Generation and Partitioning**: Custom scripts are provided for generating synthetic datasets and partitioning them, facilitating easy testing and deployment (dataset/generate_dataset.py and dataset/partition_dataset.go will generate and partition the dataset respectively).
generate_dataset.py takes the number of points to generate as a command line argument.

To run it, use the following command (from the root directory Q2/):
```
$ python3 dataset/generate_dataset.py --points <number_of_points>
```

To run partition_dataset.go, use the following command (from the root directory Q2/):
```
$ go run dataset/partition_dataset.go
```

This distributed approach allows for efficient processing of large datasets by leveraging the computational power of multiple servers, making it suitable for scenarios where the dataset is too large to be processed on a single machine efficiently.

## Input Format for Client

When running the client, you will be prompted to enter the following information:

1. Number of dimensions (D): An integer representing the dimensionality of the data points.
2. Number of nearest neighbors (K): An integer specifying how many nearest neighbors to find.
3. Query point: D space-separated float values representing the coordinates of the query point.

Example input (Same as with given CLI interface):
Enter the X and Y coordinates of the query point (separated by space): 1.0 2.0

Enter the value of k (number of nearest neighbors to find): 3

This will return the top 3 nearest neighbors to the point (1.0, 2.0) in the dataset:

Global K Nearest Neighbors:
Point: [1.22 1.45], Distance: 0.592368
Point: [1.09 0.5], Distance: 1.502698
Point: [-1.81 -2.63], Distance: 5.415995
(With data currently loaded in dataset/dataset.txt file)
