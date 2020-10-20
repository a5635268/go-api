package redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/matchstalk/go-admin-core/cache"
	"go-api/common/global"
	"go-api/tools"
	"go-api/tools/config"
)

func Setup()  {
	host := config.RedisConfig.Host
	port := config.RedisConfig.Port

	redis := &cache.Redis{
		ConnectOption: &redis.Options{
			Addr: host + ":" + tools.IntToString(port),
		},
	}

	err := redis.Connect()
	if err != nil{
		panic("redis connect err")
	}

	// 打开后，一直不关闭会不会有问题？
	// 可以考虑 golang pool redis
	global.Redis = redis
}
