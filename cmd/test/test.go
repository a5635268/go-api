package test

import (
	"github.com/spf13/cobra"
	"go-api/common/database"
	"go-api/common/global"
	"go-api/pkg/logger"
	"go-api/tools"
	"go-api/tools/config"
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

	usageStr := `starting test`
	global.Logger.Debug(usageStr)
}

func run() error {
	var x interface{}
	x = "Hello 沙河"
	v := x.(string)
	tools.Print(v)
	return nil
}
