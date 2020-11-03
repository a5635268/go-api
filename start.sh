#!/bin/bash
rm -f go-api
go build -o go-api main.go
echo "删除进程"
killall go-api #杀死运行中的go-api服务进程
echo "启动进程"
nohup ./go-api server -c=config/settings.dev.yml >> access.log 2>&1 & #后台启动服务将日志写入access.log文件
#sleep 3 #停3秒，由supervisor自动启动后再看进行
ps -aux | grep go-api #查看运行用的进程
netstat -tnlap #查看启动启动情况
