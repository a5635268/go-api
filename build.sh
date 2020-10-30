#!/bin/bash

# 如果是mac使用这个打包
# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-api main.go

# 如果是windows使用这个打包 自行测试
# SET CGO_ENABLED=0
# SET GOOS=linux
# SET GOARCH=amd64
# go build -o go-api main.go

# 如果是linux环境使用这个打包
rm -f go-api
go build -o go-api main.go
