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

// 此处进行扩展
// common/middleware/customerror.go:19


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

## 部署

### nginx配置

~~~
server {
    listen       80;
    server_name  go.example.com;

    access_log  /var/log/nginx/go.example.com.access.log  main;

    # 开启debug日志进行配置调试，正式改为warn
    error_log  /var/log/nginx/go.example.com.error.log  debug;
    
    location / {
      proxy_set_header   Host             $http_host;
      proxy_set_header   X-Real-IP        $remote_addr;
      proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;
      proxy_set_header   X-Forwarded-Proto  $scheme;
      rewrite ^/(.*)$ /$1 break;
      proxy_pass  http://172.28.3.126:8000;
    }
}
~~~

### supervisor配置

~~~
[program:go-api]
directory=/home/go/go-api
command=./start.sh
startsecs=3
stderr_logfile=/tmp/supervisord_go-api.log 
stdout_logfile=/tmp/supervisord_go-api.log  
user = root 
redirect_stderr = true
stdout_logfile_maxbytes = 20MB
stdout_logfile_backups = 20

[supervisord]
[supervisorctl]
~~~

**supervisor相关命令**

~~~
# 启动,启动之前先确认关了没
supervisord -c /etc/supervisord.conf

# 管理
supervisorctl status 
supervisorctl stop/start [process name or all]
supervisorctl shutdown
supervisorctl reload

# 配置开机启动
systemctl enable supervisord
# 验证开机启动是否成功： enabled
systemctl is-enabled supervisord

# 日志
# supervisor日志
tail -f /var/log/supervisor/supervisord.log

# 应用日志
tail -f /tmp/supervisord_go-api.log
~~~

### 部署脚本

~~~
Makefile # 编译脚本
build.sh # 交叉编译示例
start.sh # 重新编译后启动（注意：非supervisor情况）
~~~

### 临时部署（守护进程管理）

#### tmux

- [Tmux 使用教程](http://www.ruanyifeng.com/blog/2019/10/tmux.html)

~~~
yum -y install tmux
tmux new -s go-api
./go-api # 可以直接关掉窗口，进程不会段
ctrl + B & D # 退出
tmux ls # 查看
tmux attach -t go-api # 进入窗口
~~~

#### systemctl

- [Systemd 入门教程](http://www.ruanyifeng.com/blog/2016/03/systemd-tutorial-commands.html)

~~~
[Unit]
# 单元描述
Description=GO-API
# 在什么服务启动之后再执行本程序
After=mysql.service

[Service]
Type=simple
# 程序执行的目录
WorkingDirectory=/data/server/goapi/
# 启动的脚本命令
ExecStart=/data/server/goapi/goapi
# 重启条件
Restart=alway
# 几秒后重启
RestartSec=5

[Install]
WantedBy=multi-user.target
~~~

1.  创建应用配置文件 `/etc/systemd/system/go-api.service`, 内容如上;
2.  使用 `systemctl daemon-reload` 重新加载服务;
3.  执行 `systemctl start go-api` 来启动服务;
4.  最后执行 `systemctl status go-api` 来查看服务运行的状态信息;
5.  执行 `systemctl enable go-api` 将服务添加到开机启动项;
6.  注意：执行的 `go-api` 是使用文件名作为服务名;
7.  常见的命令有: `start(启动), stop(停止), restart(重启), status(查看运行状态), enable(添加到开机启动项), disable(将程序从开机启动中移除)`

#### screen

~~~
screen -S yourname -> 新建一个叫 yourname 的 session
screen -ls -> 列出当前所有的 session
screen -r yourname -> 回到 yourname 这个 session
screen -d yourname -> 远程detach某个 session
screen -d -r yourname -> 结束当前 session 并回到 yourname 这个 session
~~~

1. 使用命令 `screen -S go-api` 创建一个 session;
2. 在新终端窗口中执行 `./go-api` 即可；
3. 执行 `ctrl-a, ctrl-d` 暂时离开当前session;
4. 执行 `screen -r go-api` 返回命令窗口; 若返回不成功, 可能是该窗口被占用(Attached)了, 可以尝试使用 `screen -Dr go-api`;
5. 执行 `screen -X -S go-api quit` 结束程序;

### 容器部署 

~~~
FROM golang:alpine as builder
# 创建镜像1，用于编译go
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src/go-api
RUN go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,direct
COPY . .
RUN go env && go list && go build -v -a -o go-api .

# 二段编译： 创建mini镜像2，用于运行go
# docker build . -t  go-api
# docker run -p 9001:8000 go-api
FROM scratch
COPY --from=builder /go/src/go-api/config /config
COPY --from=builder /go/src/go-api/go-api /
#CMD ["/bin/bash"]
ENTRYPOINT ["/go-api","server","-c=config/settings.dev.yml"]
~~~

### 容器编排

~~~
cd docker-compose
docker-compse up -d
~~~

## DevOps

**通过github打包同步到目标主机**

~~~yaml
name: Github actions demo
on:
  push:
    branches:
      - master
jobs:
  build:
    name: GitHub Actions演示
    runs-on: ubuntu-latest
    steps:
      - name: environment prepare stage
        env:
          MY_VAR: Hi world! My name is
          FIRST_NAME: zhou
          MIDDLE_NAME: xiao
          LAST_NAME: gang
        run:
          echo $MY_VAR $FIRST_NAME $MIDDLE_NAME $LAST_NAME.
      - name: Set up Go 1.14 stage
        uses: actions/setup-go@v1
        id: go
        with:
          go-version: 1.14.7
      - name: Checkout stage
        uses: actions/checkout@master
      - name: build stage
        run: go build -o go-api .
      - name: Deploy stage
        uses: easingthemes/ssh-deploy@v2.1.2
        env:
          SSH_PRIVATE_KEY: ${{ secrets.SERVER_SSH_KEY }}
          ARGS: "-rltgoDzvO"
          REMOTE_HOST: ${{ secrets.REMOTE_HOST }}
          REMOTE_USER: ${{ secrets.REMOTE_USER }}
          SOURCE: "go-api"
          TARGET: ${{ secrets.REMOTE_TARGET }}
~~~
