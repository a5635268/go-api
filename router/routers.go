package router

import (
	"github.com/gin-gonic/gin"
	"go-api/app/controller"
	"go-api/app/validate"
	jwt "go-api/pkg/jwtauth"
)

var (
	routerNoCheckRole = make([]func(*gin.RouterGroup), 0)
)



func InitApiRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.RouterGroup {

	g := r.Group("")

	// 无需JWT认证
	sysNoCheckRoleRouter(g)

	// 需要认证
	sysCheckRoleRouterInit(g, authMiddleware)

	return g
}

func sysNoCheckRoleRouter(r *gin.RouterGroup) {

	// 可根据业务需求来设置接口版本
	v1 := r.Group("/v1") //
	v1.Any("/", validate.Test(), controller.Index)
	v1.Any("/test", controller.Test)

}

func sysCheckRoleRouterInit(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	r.GET("/token", authMiddleware.LoginHandler)
	v1 := r.Group("/v1").Use(authMiddleware.MiddlewareFunc())
	{
		v1.GET("/auth", controller.Auth)
	}
}
