package config

import "github.com/spf13/viper"

type Redis struct {
	Host string
	Port int
	Ttl int
}

func InitRedis(cfg *viper.Viper) *Redis {
	redis := &Redis{
		Host: cfg.GetString("host"),
		Port: cfg.GetInt("port"),
		Ttl: cfg.GetInt("ttl"),
	}
	return redis
}

var RedisConfig = new(Redis)

