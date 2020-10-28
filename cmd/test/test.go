package test

import (
	"github.com/spf13/cobra"
	"go-api/common/database"
	"go-api/common/global"
	"go-api/common/model"
	"go-api/pkg/logger"
	"go-api/pkg/redis"
	"go-api/tools"
	"go-api/tools/config"
	"time"
)

var (
	StartCmd  = &cobra.Command{
		Use:     "test",
		Short:   "用于测试语法的命令",
		Example: "go-api test",
		PreRun: func(cmd *cobra.Command, args []string) {
			 setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func setup() {
	//1. 读取配置
	configYml := "config/settings.dev.yml"
	config.Setup(configYml)
	//2. 设置日志
	logger.Setup()
	//4. 初始化redis-cache
	redis.Setup()
	//3. 初始化数据库链接
	database.Setup()


	usageStr := `starting test`
	global.Logger.Debug(usageStr)
}

func run() error {
	defer CountTime(time.Now())
	test()
	return nil
}


type Time struct {
	Id           int `json:"id" gorm:"primary_key;AUTO_INCREMENT" cache:"1000"`
	Name         string `json:"name"`
	Create_time  int `json:"create_time" gorm:"autoCreateTime"`      // 使用秒级时间戳填充创建时间
	Update_time  int `json:"update_time" gorm:"autoUpdateTime"`
	*model.Cache `json:"-"`
}

func (t Time) TableName() string {
  return "test_time"
}


func test()  {
	db := global.Eloquent
	var time []Time
	_ = db.Find(&time)
	global.Logger.Info(time)
	tools.Print(time)
}

// 运行时间计算
func CountTime(startTime time.Time)  {
	// 开始时间
	terminal := time.Since(startTime)
	global.Logger.File("debug").Debugf("运行时间 %v", terminal)
}
