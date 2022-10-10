/*
 * @Author: jlz
 * @Date: 2022-08-11 13:15:08
 * @LastEditTime: 2022-08-31 17:05:10
 * @LastEditors: jlz
 * @Description: 响应封装
 */

package common

import (
	"dnds_go/global"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 常用的状态码
const (
	successCode         = 0     // 成功
	invalidAuthCode     = 40100 // 无效的授权
	invalidArgumentCode = 40400 // 无效的参数
	logicExceptionCode  = 50100 // 逻辑异常
)

type Res struct {
	Code int
	Data interface{}
	Msg  string
}

// 响应结构体
func Response(c *gin.Context, code int, data interface{}, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}

// 成功的响应
func Succ(c *gin.Context, data interface{}, msg string) {
	Response(c, successCode, data, msg)
}

// 失败的通用响应
func Fail(c *gin.Context, code int, data interface{}, msg string) {
	Response(c, code, data, msg)
}

// validtor类型错误去除错误结构体名称，只保留字段名称
func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err

	}
	return res
}

// 通用参数校验错误
func BadParam(c *gin.Context, data interface{}, errMsg string) {
	Fail(c, ErrorBadParamCode, data, errMsg)
}

// 使用 gin自带错误校验
// 只返回第一条错误信息
func DefaultBadParam(c *gin.Context, err validator.ValidationErrors) {
	errMsgs := RemoveTopStruct(err.Translate(global.Trans))
	for _, err := range errMsgs {
		BadParam(c, nil, err)
		return
	}
}

// Beego 参数错误
func BBadParam(c *gin.Context, errMsg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": ErrorBadParamCode,
		"data": "",
		"msg":  errMsg,
	})
}
