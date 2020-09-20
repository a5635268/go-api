# go-api

go-api  基于gin的api脚手架

## 功能

### 基于cobra的多命令构建

```shell script
$ ./go-api.exe help
go-api

Usage:
  go-api [flags]
  go-api [command]

Available Commands:
  config      Get Application config info
  help        Help about any command
  server      Start API server
  test        用于测试语法的命令

Flags:
  -h, --help   help for go-api

Use "go-api [command] --help" for more information about a command.
```

### 基于`bee`的热更新

```shell script
go get github.com/beego/bee
bee run -runargs="server -c=config/settings.dev.yml"
```

### 错误处理（404）

```go
// 错误断言
panic("CustomError#" + strconv.Itoa(statusCode) + "#" + msg)

// common/middleware/customerror.go:19
// 此处进行扩展

// 404捕获
// router/init.go:39
```
