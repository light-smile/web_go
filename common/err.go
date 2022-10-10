package common

/*
	定义错误代码，自定义错误类型
*/
import (
	"dnds_go/global"
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
)

/*

ERROR(0),
 SUCCESS(1),
 REPEAT(2),   //重复
 BADPARAM(3), //参数错误
 UNAUTH(401), //没有权限
 TOKENEXPIRED(402), //token过期
 NOTLOGON(403), //未登录
 NOTFOUND(404), //未定义
 HYSTRIX(406); //熔断器返回
*/
const (
	// ErrorCode             = 0
	ErrorCommonCode       = 1
	ErrorBadParamCode     = 3
	ErrorNotFoundCode     = 404
	ErrorUnauthCode       = 401
	ErrorTokenExpiredCode = 402
	ErrorNotLoginCode     = 403
	EerrorLoginFaileCode  = 4

	ErrorMsg             = "fail"
	ErrorCommonMsg       = "合理错误"
	ErrorBadParamMsg     = "参数校验失败"
	ErrorNotFoundMsg     = "未定义"
	ErrorUnauthMsg       = "没有权限"
	ErrorTokenExpiredMsg = "token过期"
	ErrorNotLoginMsg     = "未登录"
	EerrorLoginFaileMsg  = "用户名或密码错误"
)

// 自定义错误类型，用于向客户端返回业务错误信息，通过类型区分系统错误，实现了error接口
type MyError struct {
	Msg   string      // 业务错误信息(自定义的错误消息提示)
	Err   error       // 系统错误信息(函数内部的报错信息)
	param interface{} // 接口请求参数(记录接口请求参数)
}

func (e *MyError) Error() string {
	if e.Err != nil {
		// 记录系统错误
		global.Logger.Error(e.Err.Error())
		// 返回给客户端可读错误
		return e.Msg
	}
	return e.Msg
}

// 生成一个myErr类型
func NewMyErr(msg string, err error) error {
	return &MyError{
		Msg: msg,
		Err: err,
	}
}

// 判断是否是MyErr类型
func IsMyErr(err error) bool {
	var e *MyError
	if ok := errors.As(err, &e); ok {
		return true
	}
	return false
}

// 如果错误死自定义业务错误,直接响应客户端并返回true
func MyErrHandle(c *gin.Context, err error) bool {
	if ok := IsMyErr(err); ok {
		Fail(c, ErrorCommonCode, nil, err.(*MyError).Error())
		return true
	}
	return false
}

// 自定义记录错误日志
func LogErr(params interface{}, err error) {
	if params != nil {
		j, errj := json.Marshal(params)
		if errj != nil {
			global.Logger.Error(err.Error())
		}
		global.Logger.Error(
			"errInfo:" + err.Error() + "  |" + "  params:" + string(j),
		)
		return
	}
	global.Logger.Error(
		"errInfo:" + err.Error(),
	)
}
