#!/bin/bash

# Set certificate directory
CERT_DIR=./certificate

# Clean up previous keys and certificates
rm $CERT_DIR/*.pem

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout $CERT_DIR/ca-key.pem -out $CERT_DIR/ca-cert.pem -subj "/C=FR/ST=Occitanie/L=Toulouse/O=RideShare Authority/OU=Certificate Authority/CN=rideshare-authority.com/emailAddress=ca@rideshare.com"

echo "CA's self-signed certificate"
openssl x509 -in $CERT_DIR/ca-cert.pem -noout -text

# 2. Generate server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout $CERT_DIR/server-key.pem -out $CERT_DIR/server-req.pem -subj "/C=FR/ST=Ile de France/L=Paris/O=RideShare Platform/OU=Server/CN=*.rideshare.com/emailAddress=server@rideshare.com"

# 3. Use CA's private key to sign the server's CSR and get back the signed certificate
openssl x509 -req -in $CERT_DIR/server-req.pem -days 60 -CA $CERT_DIR/ca-cert.pem -CAkey $CERT_DIR/ca-key.pem -CAcreateserial -out $CERT_DIR/server-cert.pem -extfile $CERT_DIR/server-ext.cnf

echo "Server's signed certificate"
openssl x509 -in $CERT_DIR/server-cert.pem -noout -text

# 4. Generate driver's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout $CERT_DIR/driver-key.pem -out $CERT_DIR/driver-req.pem -subj "/C=FR/ST=Alsace/L=Strasbourg/O=RideShare Platform/OU=Driver/CN=driver.rideshare.com/emailAddress=driver@rideshare.com"

# 5. Use CA's private key to sign the driver's CSR and get back the signed certificate
openssl x509 -req -in $CERT_DIR/driver-req.pem -days 60 -CA $CERT_DIR/ca-cert.pem -CAkey $CERT_DIR/ca-key.pem -CAcreateserial -out $CERT_DIR/driver-cert.pem -extfile $CERT_DIR/client-ext.cnf

echo "Driver's signed certificate"
openssl x509 -in $CERT_DIR/driver-cert.pem -noout -text

# 6. Generate rider's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout $CERT_DIR/rider-key.pem -out $CERT_DIR/rider-req.pem -subj "/C=FR/ST=Alsace/L=Strasbourg/O=RideShare Platform/OU=Rider/CN=rider.rideshare.com/emailAddress=rider@rideshare.com"

# 7. Use CA's private key to sign the rider's CSR and get back the signed certificate
openssl x509 -req -in $CERT_DIR/rider-req.pem -days 60 -CA $CERT_DIR/ca-cert.pem -CAkey $CERT_DIR/ca-key.pem -CAcreateserial -out $CERT_DIR/rider-cert.pem -extfile $CERT_DIR/client-ext.cnf

echo "Rider's signed certificate"
openssl x509 -in $CERT_DIR/rider-cert.pem -noout -text

# Optional: Display the contents of the certificates (server, driver, and rider)
echo "Displaying certificate details:"
openssl x509 -in $CERT_DIR/server-cert.pem -noout -text
openssl x509 -in $CERT_DIR/driver-cert.pem -noout -text
openssl x509 -in $CERT_DIR/rider-cert.pem -noout -text
