#!/bin/bash

set -e

# 依赖的组件 protoc.exe protoc-gen-go.exe protoc-gen-go-grpc.exe
# 将依赖组件拷贝到 GOPATH/bin 路径下执行脚本
protoc --go_out=. --go-grpc_out=require_unimplemented_servers=false:. connect.proto

mv ./pb/*.go ../pb/

rm -rf pb