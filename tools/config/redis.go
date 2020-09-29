package config

import "github.com/spf13/viper"

type Redis struct {
	Host string
	Port int
}

func InitRedis(cfg *viper.Viper) *Redis {
	redis := &Redis{
		Host: cfg.GetString("host"),
		Port: cfg.GetInt("port"),
	}
	return redis
}

var RedisConfig = new(Redis)

