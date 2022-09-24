#!/bin/bash
set -ex
cd src
export PATH=$PATH:/usr/local/go/bin
go version
go test -v
go build -o go-service main.go mysql.go
ls -lh
