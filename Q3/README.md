# Homework-4: gRPC

```
Names : Mitansh Kayathwal, Pradeep Mishra
Roll Nos: 2021101026, 2023801013
Branch : CSE, PhD
Course : Distributed Systems, Monsoon '24
```

# **_Directory Structure_**

```
📁 Q3
├── 📁 certificate
│   └── 📄 gen_certificate.sh
│   └── 📄 server-ext.cnf
│   └── 📄 client-ext.cnf
├── 📁 client
│   └── 📄 driver.go
│   └── 📄 rider.go
├── 📁 protofiles
│   └── 📄 ride_grpc.pb.go
│   └── 📄 ride.pb.go
│   └── 📄 ride.proto
├── 📁 server
│   └── 📄 default.etcd
│   └── 📄 main.go
├── 📄 go.mod
├── 📄 go.sum
├── 📄 Makefile
├── 📄 README.md
```

## Running the Code

To run the MyUber Application, follow these steps (from the root directory Q3/):

1. Build the proto-files:
   ```
   make proto
   ```
   
   This script will clean up previous builds and generate the necessary proto files.

2. Generate the certificates required for mTLS authorization:
   ```
   chmod +x ./certificate/gen_certificate.sh
   make certificate
   ```

3. Start the etcd client (You would need to install this):
    ```
    $ etcd
    ```

    More documentation for installation can be found at:
    - https://pkg.go.dev/go.etcd.io/etcd/client/v3
    - https://etcd.io/docs/v3.2/install/


2. Start multiple server instances in different terminals (upto 3 servers support provided in Makefile):
   ```
   $ make server1
   ```

   ```
   $ make server2
   ```

   ```
   $ make server3
   ```

3. In a separate terminal, start instances of the 2 clients:

   - To start rider client, you can use the following according to your load balancing requirement:
        ```
        $ make rider_roundrobin
        ```

        ```
        $ make rider_pickfirst
        ```

    - To start driver client, you can use the following:
        ```
        $ make driver
        ```

5. Follow the prompts in the clients rider and driver to interact with the MyUber system. You can test the application as per your wish. All requirements for the system have been fulfilled. From the rider client, you get an option to request a ride, get it's status (if a ride is ongoing) and exit the application. From the driver end you can accept/reject ride requests.