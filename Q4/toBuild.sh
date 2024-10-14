#!/bin/bash

ORANGE='\033[0;33m'
GREEN='\033[0;32m'
NC='\033[0m'

# Running make clean to remove any previous build artifacts
echo -e "${ORANGE}Cleaning up previous builds...${NC}"
make clean

# Running make proto to generate proto files
echo -e "${ORANGE}Generating proto files...${NC}"
make proto

echo -e "${GREEN}Build completed.${NC}"