package global

import (
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/os/glog"
	"github.com/matchstalk/go-admin-core/cache"
	"github.com/robfig/cron/v3"
	"go-api/common/config"
	"go-api/pkg/queue/redis"
	"gorm.io/gorm"
)

const (
	Version = "0.0.2"
)

var Cfg config.Conf = config.DefaultConfig()

var GinEngine *gin.Engine

var Eloquent *gorm.DB

var GADMCron *cron.Cron

var Redis *cache.Redis

var Queue1 *redis.Queue

var Queue2 *redis.Queue

var (
	Source string
	Driver string
	DBName string
)

var (
	Logger        *glog.Logger
	JobLogger     *glog.Logger
	RequestLogger *glog.Logger
)
