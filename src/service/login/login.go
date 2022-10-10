package loginService

import (
	"dnds_go/common"
	"dnds_go/global"
	loginDao "dnds_go/src/dao/login"
	loginModel "dnds_go/src/models/login"
	usersModel "dnds_go/src/models/users"
	"dnds_go/tool"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var user loginModel.LoginReq
	if ok := common.BindJson(c, &user); !ok {
		return
	}
	// 进行登录叫校验，通过获取token信息
	res, err := loginDao.Login(user)
	if err != nil {
		if ok := common.IsMyErr(err); ok {
			common.Fail(c, common.EerrorLoginFaileCode, nil, err.Error())
			return
		}
		global.Logger.Error(err.Error())
		common.Fail(c, common.EerrorLoginFaileCode, nil, "系统错误")
		return
	}
	// 生成uuid，作为客户端唯一标识，如果请求未携带uuid视为离线
	uuid := tool.NewUuid()
	// 将该登录用户名的uuid更新
	if err := loginDao.UpdateUuid(user, uuid); err != nil {
		global.Logger.Error(err.Error())
		common.Fail(c, http.StatusOK, nil, "登录失败")
		return
	}
	common.Succ(c, res, "登录成功")
}

func CreateUser(c *gin.Context) {
	var u usersModel.User
	if ok := common.BindJson(c, &u); !ok {
		return
	}
	user := usersModel.User{
		UserName:     u.UserName,
		UserLogin:    u.UserLogin,
		UserPassword: u.UserPassword,
	}

	if err := user.Create(); err != nil {
		// 判断是否是账号已存在
		if isMyErr := common.IsMyErr(err); isMyErr {
			common.Fail(c, 1, nil, err.Error())
			return
		}
		// 如果是数据库报错，记录错误，返回创建失败
		global.Logger.Error(err.Error())
		common.Fail(c, 1, nil, "用户创建失败")
		return
	}

	common.Succ(c, nil, "用户创建成功")
}
