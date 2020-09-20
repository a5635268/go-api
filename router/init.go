package router

import (
	"github.com/gin-gonic/gin"
	"go-api/common/global"
	"go-api/common/middleware"
	"go-api/common/middleware/handler"
	"go-api/tools"
	config2 "go-api/tools/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitRouter(r *gin.Engine) *gin.Engine {
	if config2.SslConfig.Enable {
		r.Use(handler.TlsHandler())
	}

	// 加载mysql数据库
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: global.Cfg.GetDb().DB}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		global.Logger.Fatal(tools.Red("mysql connect error :"), err)
	}
	gormDB := make(map[string]*gorm.DB)
	gormDB["*"] = db

	r.Use(middleware.WithContextDb(gormDB))
	middleware.InitMiddleware(r)

	// 注册业务路由
	InitApiRouter(r)

	// 404捕获
	r.NoRoute(handler.HandleNotFound())

	return r
}
