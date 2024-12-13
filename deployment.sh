#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

graceful_exit() {
  echo -e "${YELLOW}Exiting gracefully...${NC}"
  exit 0
}

trap graceful_exit SIGINT SIGTERM

check_file_exists() {
  if [ ! -f "$1" ]; then
    echo -e "${RED}Error: File '$1' not found!${NC}"
    exit 0
  fi
}

check_command() {
  if ! command -v "$1" &> /dev/null; then
    echo -e "${RED}Error: $1 is not installed. Please install it first.${NC}"
    exit 0
  fi
}

check_command "docker"

check_file_exists ".env"

LOGFILE="deploy_$(date +%Y%m%d_%H%M%S).log"
exec > >(tee -a "$LOGFILE") 2>&1

START_TIME=$(date +%s)

echo -e "${YELLOW}Removing old Docker image...${NC}"
docker image rm -f ervis

echo -e "${YELLOW}Building new Docker image...${NC}"
docker build -t ervis .

echo -e "${YELLOW}Bringing down existing Docker containers...${NC}"
docker compose down

echo -e "${YELLOW}Starting up Docker containers...${NC}"
docker compose up -d

END_TIME=$(date +%s)
EXECUTION_TIME=$((END_TIME - START_TIME))

echo -e "${GREEN}Deployment started successfully!${NC}"
echo -e "${GREEN}Total execution time: ${EXECUTION_TIME} seconds.${NC}"

echo -e "${YELLOW}Logs saved to ${LOGFILE}${NC}"

echo -e "${YELLOW}Press any key to exit...${NC}"
read -n 1 -s

graceful_exit