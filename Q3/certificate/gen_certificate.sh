#!/bin/bash

# Clean up previous keys and certificates
rm *.pem

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=FR/ST=Occitanie/L=Toulouse/O=RideShare Authority/OU=Certificate Authority/CN=rideshare-authority.com/emailAddress=ca@rideshare.com"

echo "CA's self-signed certificate"
openssl x509 -in ca-cert.pem -noout -text

# 2. Generate server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=FR/ST=Ile de France/L=Paris/O=RideShare Platform/OU=Server/CN=*.rideshare.com/emailAddress=server@rideshare.com"

# 3. Use CA's private key to sign the server's CSR and get back the signed certificate
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf

echo "Server's signed certificate"
openssl x509 -in server-cert.pem -noout -text

# 4. Generate driver's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout driver-key.pem -out driver-req.pem -subj "/C=FR/ST=Alsace/L=Strasbourg/O=RideShare Platform/OU=Driver/CN=driver.rideshare.com/emailAddress=driver@rideshare.com"

# 5. Use CA's private key to sign the driver's CSR and get back the signed certificate
openssl x509 -req -in driver-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out driver-cert.pem -extfile client-ext.cnf

echo "Driver's signed certificate"
openssl x509 -in driver-cert.pem -noout -text

# 6. Generate rider's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout rider-key.pem -out rider-req.pem -subj "/C=FR/ST=Alsace/L=Strasbourg/O=RideShare Platform/OU=Rider/CN=rider.rideshare.com/emailAddress=rider@rideshare.com"

# 7. Use CA's private key to sign the rider's CSR and get back the signed certificate
openssl x509 -req -in rider-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out rider-cert.pem -extfile client-ext.cnf

echo "Rider's signed certificate"
openssl x509 -in rider-cert.pem -noout -text

# Optional: Display the contents of the certificates (server, driver, and rider)
echo "Displaying certificate details:"
openssl x509 -in server-cert.pem -noout -text
openssl x509 -in driver-cert.pem -noout -text
openssl x509 -in rider-cert.pem -noout -text
