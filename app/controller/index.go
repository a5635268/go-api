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
	args :=  c.DefaultQuery("args", "1")
	if args == "2"{
		panic("cuowu")
	}
	tools.Print(result,err)
	app.OK(c, result, "hahaha")
}

func Auth(c *gin.Context)  {
	uid,_ :=  c.Get("uid")
	username,_ := c.Get("username")
	openid,_ := c.Get("openid")
	data,_ := c.Get("JWT_PAYLOAD")
	result := map[string]interface{}{
		"uid" : uid,
		"username" : username,
		"openid" : openid,
		"data" : data,
	}
	app.OK(c, result, "")
}
