#!/bin/bash
set -ex
cd src
go test -v
go build -o go-service
ls -lh
