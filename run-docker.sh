#!/bin/bash
set +e
docker rm go-service
set -e
make linux
docker build -t go-service .
docker run -t -p 8081:8081 --name go-service -it go-service
