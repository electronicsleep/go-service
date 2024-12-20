#!/bin/bash
UNAME=$(uname -m)
echo "UNAME: $UNAME"
# BUILD=arm64
BUILD=amd64
REPO=localhost:32000

IMG=$(docker images -q "$REPO/go-service")
if [[ ! -z "$IMG" ]]; then
  docker rmi --force $IMG
fi
set -e

# build
export GOOS=linux
cd src && go build -o go-service .
cd -

# publish
docker build --platform linux/$BUILD . -t $REPO/go-service:latest
docker push $REPO/go-service:latest
