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

### 基于`bee`和`gin`的热更新

```shell script
go get github.com/beego/bee
bee run -runargs="server -c=config/settings.dev.yml"

# 一般用于脚本调试用
go get github.com/codegangsta/gin
gin run main.go
```

### 错误处理（404捕获）

```go
// 错误断言
panic("CustomError#" + strconv.Itoa(statusCode) + "#" + msg)

// common/middleware/customerror.go:19
// 此处进行扩展

// 404捕获
// router/init.go:39
```

### 日志使用

```go
// 详情： https://www.bookstack.cn/read/GoFrame-1.13/os-glog-chain.md
// 带f的就是字符占位，用法和fmt.printf一样。
global.Logger.Info // Warning,Debug,Notice,Error,Fatal(直接断掉退出程序),Panic(写入日志，并且抛出错误)
global.Logger.Infof

// 分类
global.Logger.Cat("debug") // File

// pkg/logger/logger.go:14 日志配置

// 提供map/struct类型参数  打印json
```

### 参数验证

```go
// app/validate
// 参考：https://github.com/gookit/validate/blob/master/README.zh-CN.md

// 常用自定义验证器写在：app/validate/init.go:11
// 中间件引入验证器
v1.Any("/", validate.Test(), controller.Index)
```

### redis支持

```go
// 普通命令
res,_ := global.Redis.Get("zhouzhou")
tools.Print(res)

// redis原生命令支持
zsetKey := "zsettest"
redis := global.Redis.GetClient()
ret, err := redis.ZRevRangeWithScores(zsetKey, 0, 2).Result()
if err != nil {
    fmt.Printf("zrevrange failed, err:%v\n", err)
    return
}
for _, z := range ret {
    fmt.Println(z.Member, z.Score)
}

// redis队列 queredis "go-api/pkg/queue/redis"
message := queue.Msg{100, "队列1", 43}
msg := queredis.NewMessage("", time.Now().Add(time.Second*1), message)
_, err := global.Queue1.Publish(msg)
tools.HasError(err,"")
```

### Gorm定制化

```go

```
