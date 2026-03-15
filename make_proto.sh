#!/usr/bin/env bash

protoc -I . --go_out=./proto --go-grpc_out=./proto --go_opt=module=github.com/Compogo/db-client/proto \
 ./db-client/proto/*.proto
