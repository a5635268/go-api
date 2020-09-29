package queue

import (
	"context"
	"fmt"
	"go-api/common/global"
	queue "go-api/pkg/queue/redis"
)

type Msg struct {
	MsgId   int    `json:"msg_id"`
	MsgBody string `json:"body"`
	UserId  int    `json:"uid"`
}

func handerFuncQ1(msg queue.Message) {
	fmt.Println("queue1 消费一条消息：=========")
	fmt.Printf("%#v\n", msg)

	// 转map
	m := msg.Body.(map[string]interface{})
	fmt.Println(m["msg_id"], m["body"], m["uid"])
}

func handerFuncQ2(msg queue.Message) {
	fmt.Println("queue2 消费一条消息：=========")
	fmt.Printf("%#v\n", msg)

	// 转map
	m := msg.Body.(map[string]interface{})
	fmt.Println(m["msg_id"], m["body"], m["uid"])
}

func SetUp(){
	// 启动队列1
	queue1 := queue.NewQueue(
		context.TODO(),
		global.Redis.GetClient(),
		queue.WithTopic("queue1"),
		queue.WithHandler(handerFuncQ1))
	global.Queue1 = queue1
	global.Queue1.Start()
	global.Logger.Info("队列1启动成功")

	// 启动队列2
	queue2 := queue.NewQueue(
		context.TODO(),
		global.Redis.GetClient(),
		queue.WithTopic("queue2"),
		queue.WithHandler(handerFuncQ2))
	global.Queue2 = queue2
	global.Queue2.Start()
	global.Logger.Info("队列2启动成功")
}

