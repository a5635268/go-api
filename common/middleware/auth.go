package middleware

import (
	"time"

	"go-api/common/middleware/handler"
	jwt "go-api/pkg/jwtauth"
	"go-api/tools/config"
)

func AuthInit() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "test zone",
		Key:             []byte(config.ApplicationConfig.JwtSecret),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		// 数据加载
		PayloadFunc:     handler.PayloadFunc,
		// 获取数据
		IdentityHandler: handler.IdentityHandler,
		// 数据验证: 登录
		Authenticator:   handler.Authenticator,
		//中间件向控制器传值用
		Authorizator:    handler.Authorizator,
		// 未授权时调用
		Unauthorized:    handler.Unauthorized,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})

}
