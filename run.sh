#!/bin/bash
set -e
export writerDatasource=127.0.0.1
export readerDatasource=127.0.0.1
export datasourcePassword=password
cd src
go test -v
go build main.go mysql.go
./main
