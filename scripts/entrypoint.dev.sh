#!/bin/bash
set -e

echo "migrating dev db..."
go run cmd/dbmigrate/main.go

echo "downloading CompileDaemon..."
# disable go modules to avoid this package from getting into go.mod
# as we only need it locally to watch and rebuild server on change
GO111MODULE=off go get github.com/githubnemo/CompileDaemon

echo "starting test script deamon..."
nohup CompileDaemon --build="go build -o testserver cmd/testserver/main.go" --command=./testserver &

echo "starting main deamon..."
CompileDaemon --build="go build -o main cmd/burst/main.go" --command=./main

