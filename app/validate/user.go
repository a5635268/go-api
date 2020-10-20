package validate

import (
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
	"go-api/tools"
)

// 参考：https://github.com/gookit/validate/blob/master/README.zh-CN.md
type User struct {
	Name string `json:"name" validate:"myCheck1|customValidator|required" message:"customValidator:自定义错误|required:名称不为空|myCheck1:自我怀疑"`
	BillingAddressId uint `json:"billing_address_id" validate:"required"`
	ShippingAddressId uint `json:"shipping_address_id"`
}


// CustomValidator 定义在结构体中的自定义验证器
func (f User) CustomValidator(val string) bool {
	return  true
}

// Translates 你可以自定义字段翻译
func (f User) Translates() map[string]string {
	return validate.MS{
		"Name": "用户名称",
	}
}

// 优先级高
// Messages 您可以自定义验证器错误消息
func (f User) Messages() map[string]string {
	return validate.MS{
		"required": "{field} 不能为空",
	}
}

func Test() gin.HandlerFunc{
	return func(c *gin.Context) {
		var data User
		err := c.ShouldBind(&data)
		tools.HasError(err, "", VerifyErrorCode)

		v := validate.Struct(data)
		tools.Assert(v.Validate(),v.Errors.One(),VerifyErrorCode)

		c.Set(c.Request.Method + "Param",data)
		c.Next()
	}
}
