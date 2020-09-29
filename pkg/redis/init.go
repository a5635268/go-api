package redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/matchstalk/go-admin-core/cache"
	"github.com/matchstalk/redisqueue"
	"go-api/common/global"
	"go-api/tools"
	"go-api/tools/config"
	"time"
)

func Setup()  {
	host := config.RedisConfig.Host
	port := config.RedisConfig.Port

	redis := &cache.Redis{
		ConnectOption: &redis.Options{
			Addr: host + ":" + tools.IntToString(port),
		},
		// 只有redis5版本以上才支持
		ConsumerOptions: &redisqueue.ConsumerOptions{
			VisibilityTimeout: 60 * time.Second,
			BlockingTimeout:   5 * time.Second,
			ReclaimInterval:   1 * time.Second,
			BufferSize:        100,
			Concurrency:       10,
		},
		ProducerOptions: &redisqueue.ProducerOptions{
			StreamMaxLength:      100,
			ApproximateMaxLength: true,
		},
	}

	err := redis.Connect()
	if err != nil{
		panic("redis connect err")
	}
	global.Redis = redis
}
