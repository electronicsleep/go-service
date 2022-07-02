#!/bin/bash
set -e
# go test -v
export writerDatasource=127.0.0.1
export readerDatasource=127.0.0.1
export datasourcePassword=password
go build main.go mysql.go
./main
