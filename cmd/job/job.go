package job

import (
	"github.com/spf13/cobra"
	"go-api/common/global"
	"go-api/pkg/cronjob"
	"go-api/pkg/logger"
	"go-api/tools/config"
)

var (
	StartCmd = &cobra.Command{
		Use:     "job",
		Short:   "用于开启定时任务的脚本",
		Example: "go-api job",
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

/*	//4. 初始化redis-cache
	redis.Setup()

	//3. 初始化数据库链接
	database.Setup()*/

	usageStr := `starting job`
	global.Logger.Debug(usageStr)
}

func run() error {
	cronjob.InitJob()
	cronjob.Setup()
	return nil
}
