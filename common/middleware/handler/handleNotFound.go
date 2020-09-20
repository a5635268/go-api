package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-api/tools"
)

func HandleNotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := c.Request.Method + " " + c.Request.URL.String()
		tools.HasError(errors.New(request + " not found"),"")
		return
	}
}
