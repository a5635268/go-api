package api

import (
	"context"
	"fmt"
	"go-api/pkg/queue"
	"go-api/pkg/redis"
	"go-api/router"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go-api/common/database"
	"go-api/common/global"
	"go-api/pkg/logger"
	"go-api/tools"
	"go-api/tools/config"
)

var (
	configYml string
	port      string
	mode      string
	StartCmd  = &cobra.Command{
		Use:          "server",
		Short:        "Start API server",
		Example:      "go-api server -c config/settings.yml",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

var echoTimes int

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.dev.yml", "Start server with dev configuration file")
	StartCmd.PersistentFlags().StringVarP(&port, "port", "p", "8000", "Tcp port server listening on")
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,test,prod")
}

func setup() {
	//1. 读取配置
	config.Setup(configYml)
	//2. 设置日志
	logger.Setup()
	//3. 初始化数据库链接
	database.Setup()
	//4. 初始化redis-cache
	redis.Setup()
	//5. 初始化queue
	queue.SetUp()

	usageStr := `starting api server`
	global.Logger.Info(usageStr)
}

func run() error {
	if viper.GetString("settings.application.mode") == string(tools.ModeProd) {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := global.Cfg.GetEngine()
	if engine == nil {
		engine = gin.New()
	}
	r := router.InitRouter(engine.(*gin.Engine))
	global.Cfg.SetEngine(r)
	srv := &http.Server{
		Addr:    config.ApplicationConfig.Host + ":" + config.ApplicationConfig.Port,
		Handler: global.Cfg.GetEngine(),
	}
	go func() {
		// 服务连接
		if config.SslConfig.Enable {
			if err := srv.ListenAndServeTLS(config.SslConfig.Pem, config.SslConfig.KeyStr); err != nil && err != http.ErrServerClosed {
				global.Logger.Fatal("listen: ", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				global.Logger.Fatal("listen: ", err)
			}
		}
	}()
	content, _ := ioutil.ReadFile("./static/go-api.txt")
	fmt.Println(tools.Red(string(content)))
	tip()
	fmt.Println(tools.Green("Server run at:"))
	fmt.Printf("-  Local:   http://localhost:%s/ \r\n", config.ApplicationConfig.Port)
	fmt.Printf("-  Network: http://%s:%s/ \r\n", tools.GetLocaHonst(), config.ApplicationConfig.Port)
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	fmt.Printf("%s Shutdown Server ... \r\n", tools.GetCurrentTimeStr())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.Logger.Fatal("Server Shutdown:", err)
	}
	global.Logger.Println("Server exiting")
	return nil
}

func tip() {
	usageStr := `欢迎使用 ` + tools.Green(`go-api `+global.Version) + ` 可以使用 ` + tools.Red(`-h`) + ` 查看命令`
	fmt.Printf("%s \n\n", usageStr)
}
