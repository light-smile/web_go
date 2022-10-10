package loginDao

import (
	"dnds_go/common"
	"dnds_go/global"
	loginModel "dnds_go/src/models/login"
	usersModel "dnds_go/src/models/users"
	"dnds_go/tool"
)

func Login(u loginModel.LoginReq) (loginModel.LoginRes, error) {
	var user usersModel.User
	var res loginModel.LoginRes
	global.L.Lock()
	find := global.DB.Where("user_login = ?", u.UserLogin).Find(&user)
	if find.Error != nil {
		return res, find.Error
	}
	if find.RowsAffected == 0 {
		return res, common.NewMyErr("用户名不存在", nil)
	}
	global.L.Unlock()
	upsM := tool.MD5(u.UserPassWord)
	if upsM != user.UserPassword {
		return res, common.NewMyErr("密码错误", nil)
	}
	accessToken := common.GetLoginToken(u)
	res = loginModel.LoginRes{
		AccessToken: accessToken,
	}
	return res, nil
}

func UpdateUuid(u loginModel.LoginReq, uuid string) (err error) {
	global.L.Lock()
	defer global.L.Unlock()
	res := global.DB.Table("sys_user").Where("user_login = ?", u.UserLogin).Update("client_ident", uuid)
	err = res.Error
	return
}
