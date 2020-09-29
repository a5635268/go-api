package controller

import (
	"github.com/gin-gonic/gin"
	"go-api/tools"
	"go-api/tools/app"
)

func Index(c *gin.Context) {
	result, err := c.Get("_data")
	tools.Print(result,err)
	app.OK(c, result, "")
}

func Test(c *gin.Context) {
	result, err := c.Get("_data")
	tools.Print(result,err)
	app.OK(c, result, "")
}
