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

# 命令扩展：cmd/cobra.go:42
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

#### 缓存定制

```go
type Time struct {
    // 主键tag：cache。缓存时间。默认取config.redis.ttl
	Id int `json:"id" gorm:"primary_key;AUTO_INCREMENT" cache:"1000"`
	Name string `json:"name"`
	CreateTime int `json:"create_time" gorm:"autoCreateTime"` // 使用秒级时间戳填充创建时间      
	UpdateTime int `json:"update_time" gorm:"autoUpdateTime"`
    // 嵌入匿名结构体，开启缓存
	models.Cache `json:"-"`
}
db := global.Eloquent
var time []Time
_ = db.Find(&time)
```

**Gorm回调笔记**

```go
// 当前 Model 的所有字段
db.Statement.Schema.Fields

// 当前 Model 的所有主键字段
db.Statement.Schema.PrimaryFields

// 优先主键字段：带有 db 名为 `id` 或定义的第一个主键字段。
db.Statement.Schema.PrioritizedPrimaryField

// 当前 Model 的所有关系
db.Statement.Schema.Relationships

// 当前实例化的model(取值结果集)
db.Statement.ReflectValue

// 根据 db 名或字段名查找字段
field := db.Statement.Schema.LookUpField("Name")
```

### jwt支持

~~~go
// 登录鉴权代码： common/middleware/handler/auth.go:35

// token获取途径，header: Authorization, query: token, cookie: jwt

// 控制器获取
uid,_ :=  c.Get("uid")
username,_ := c.Get("username")
openid,_ := c.Get("openid")
data,_ := c.Get("JWT_PAYLOAD")
result := map[string]interface{}{
    "uid" : uid,
    "username" : username,
    "openid" : openid,
    "data" : data,
}
app.OK(c, result, "")
~~~

### 计划任务

~~~go
// 配置计划任务： pkg/cronjob/joblist.go:8
{
    InvokeTarget:   "ExamplesOne",
    Name:           "exec test",
    JobId:          1,
    EntryId:        0, 
    CronExpression: "*/5 * * * * *", // 以秒开始，其它参考Linux crontab
    Args:           "参数",
    JobType:           2,  // 1，http。2，exec函数调用。
}

// 启动计划任务
// go build mian.go && ./go-api job

// 停掉计划任务： 杀掉进程就可以了
~~~
