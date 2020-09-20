package controller

import (
	"github.com/gin-gonic/gin"
	"go-api/tools/app"
)

func Index(c *gin.Context) {
	result := "默认控制器"
	app.OK(c, result, "")
}
