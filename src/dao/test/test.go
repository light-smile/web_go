package testDao

import (
	"dnds_go/global"
	testModel "dnds_go/src/models/test"

	"go.uber.org/zap"
)

// FindOrCreate 创建用户表 如果有同名用户直接返回该条用户，没有则创建
func FindOrCreate(user testModel.User) ([]testModel.User, error) {
	var userres []testModel.User
	if err := user.FindByName(&userres); err != nil {
		global.Logger.Error("user findByName err:"+err.Error(),
			zap.Any("params:", user),
		)
		return nil, err
	}
	if len(userres) != 0 {
		return userres, nil
	}
	if err := user.Create(); err != nil {
		global.Logger.Error("user create err:"+err.Error(),
			zap.Any("params:", user),
		)
		return nil, err
	}
	if err := user.FindByName(&userres); err != nil {
		global.Logger.Error("user findByName err:"+err.Error(),
			zap.Any("params:", user),
		)
		return nil, err
	}
	return userres, nil
}

func Create() ([]testModel.User, error) {
	user := testModel.User{
		Name: "测试性能",
		Age:  00000,
	}
	//
	global.L.Lock()
	if err := user.Create(); err != nil {
		return nil, err
	}
	global.L.Unlock()
	var result []testModel.User
	global.L.Lock()
	if err := user.FindByName(&result); err != nil {
		return nil, err
	}
	global.L.Unlock()
	return result, nil
}
