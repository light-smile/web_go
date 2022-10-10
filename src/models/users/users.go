package usersModel

import (
	"dnds_go/common"
	"dnds_go/global"
	"dnds_go/tool"
	"time"
)

type User struct {
	ID           int    `gorm:"type:int;column:ID;primaryKey;autoIncrement;not null"`
	UserName     string `json:"userName" binding:"required" gorm:"type:varchar;column:USER_NAME;not null;"`
	UserLogin    string `json:"userLogin" binding:"required" gorm:"type:varchar;column:USER_LOGIN;not null;unique"`
	UserPassword string `json:"userPassword"  gorm:"type:varchar;column:USER_PASSWORD;not null"`
	Status       int    `gorm:"type:int;column:STATUS;" `
	CreateTime   string `gorm:"type:varchar;column:CREATE_TIME;"`
	Remark       string `gorm:"type:varchar;column:REMARK"`
	UserLevel    *uint  `gorm:"type:tinyint;column:USER_LEVEL;default:1"`
}

func (u *User) TableName() string {
	return "SYS_USER"
}

func (u *User) Create() error {
	var exit User
	global.L.Lock()
	defer global.L.Unlock()
	res := global.DB.Where("user_login = ?", u.UserLogin).Find(&exit)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected > 0 {
		return common.NewMyErr("用户已存在", nil)
	}
	// 对密码进行md5 加密
	u.UserPassword = tool.MD5(u.UserPassword)
	u.CreateTime = time.Now().Format(global.TimeFormatYmdhis)
	if err := global.DB.Create(u).Error; err != nil {
		return err
	}
	return nil
}
