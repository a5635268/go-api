package validate

import "github.com/gookit/validate"

var (
	// 验证器错误
	ErrorCode = 1011
)


func init(){
	m := map[string]interface{}{
		"myCheck1": func(val interface{}) bool {
			return true
		},
		"mobile" : func(val interface{}) bool {
			return true
		},
	}
	// 此处做validate的全局配置，比如全局验证器
	validate.AddValidators(m)
}
