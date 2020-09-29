package router

import (
	"github.com/gin-gonic/gin"
	"go-api/app/controller"
	"go-api/app/validate"
)

var (
	routerNoCheckRole = make([]func(*gin.RouterGroup), 0)
)

func InitApiRouter(r *gin.Engine) *gin.Engine {
	// 可根据业务需求来设置接口版本
	v1 := r.Group("/v1") //
	v1.Any("/", validate.Test(), controller.Index)
	v1.Any("/test", controller.Test)
	return r
}
