package test

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-api/common/database"
	"go-api/common/global"
	"go-api/pkg/logger"
	"go-api/pkg/queue"
	"go-api/pkg/redis"
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
	//3. 初始化数据库链接
	database.Setup()
	//4. 初始化redis-cache
	redis.Setup()
	//5. 初始化queue
	queue.SetUp()

	usageStr := `starting test`
	global.Logger.Debug(usageStr)
}

func run() error {
	defer CountTime(time.Now())
	test()
	return nil
}

func test()  {

}

func test1() {
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
}

// 运行时间计算
func CountTime(startTime time.Time)  {
	// 开始时间
	terminal := time.Since(startTime)
	global.Logger.File("debug").Debugf("运行时间 %v", terminal)
}
