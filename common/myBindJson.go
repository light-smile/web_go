package common

import (
	"dnds_go/global"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 自定义绑定json数据，并对请求参数做校验 校验通过返回true
func BindJson(c *gin.Context, st interface{}) bool {
	if err := c.ShouldBindJSON(st); err != nil {
		// 判断是否通过校验
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 系统错误，记录日志
			global.Logger.Error(err.Error())
			BadParam(c, nil, "参数格式有误")
			return false
		}
		// 参数校验错误
		DefaultBadParam(c, errs)
		return false
	}
	return true
}
