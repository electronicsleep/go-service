#!/bin/bash
set -ex
cd src
go test -v
go build -o go-service main.go mysql.go
ls -lh
