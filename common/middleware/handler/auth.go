package handler

import (
	"github.com/gin-gonic/gin"
	"go-api/app/model"
	"go-api/common/global"
	jwt "go-api/pkg/jwtauth"
	"go-api/tools"
	"net/http"
)

func PayloadFunc(data interface{}) jwt.MapClaims {
	if user, ok := data.(map[string]interface{}); ok {
		return jwt.MapClaims{
			jwt.IdentityKey: user["uid"],
			jwt.UserNameKey: user["nickname"],
			jwt.OpenIdKey:   user["openid"],
		}
	}
	return jwt.MapClaims{}
}

func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return map[string]interface{}{
		jwt.IdentityKey: claims["identity"],
		jwt.UserNameKey: claims["nickname"],
		jwt.Uid:         claims["identity"],
		jwt.OpenIdKey:   claims["openid"],
	}
}

// 微信用户登录就写在这里
func Authenticator(c *gin.Context) (interface{}, error) {
	// 伪代码
	var user model.SysUser
	res := global.Eloquent.First(&user)
	if res.Error == nil {
		data, _ := tools.StructToMap(user)
		return data, nil
	}
	return nil, jwt.ErrFailedAuthentication
}

func Authorizator(data interface{}, c *gin.Context) bool {
	if user, ok := data.(map[string]interface{}); ok {
		c.Set("uid", user[jwt.Uid])
		c.Set("username", user[jwt.UserNameKey])
		c.Set("openid", user[jwt.OpenIdKey])
		return true
	}
	return false
}

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  message,
	})
}
