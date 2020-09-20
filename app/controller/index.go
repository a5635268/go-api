package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-api/tools"
	"go-api/tools/app"
)

func Index(c *gin.Context) {
	result := "默认控制器222333"
	err :=  errors.New("您输入不是双数")
	tools.HasError(err,"")
	app.OK(c, result, "")
}
