package router

import (
	"github.com/gin-gonic/gin"
	"go-api/app/controller"
)

func InitApiRouter(r *gin.Engine) *gin.Engine {

	// 可根据业务需求来设置接口版本
	v1 := r.Group("/v1")
	v1.GET("/", controller.Index)

	return r
}
