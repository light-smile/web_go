/*
 * @Author: jlz
 * @Date: 2022-08-10 16:27:25
 * @LastEditTime: 2022-08-31 15:50:01
 * @LastEditors: jlz
 * @Description:
 */

package provider

import (
	"dnds_go/config"
	"dnds_go/global"
	constantModel "dnds_go/src/models/constant"
	gatewayModel "dnds_go/src/models/gateway"
	usersModel "dnds_go/src/models/users"
	"log"

	// "gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/gorm"
)

// 初始化，数据库连接
func InitDatabase() {
	global.DB = gormSqlite()
	// db.AutoMigrate(&User{})
	// user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
	// result := db.Create(&user)
	// if result.Error != nil {
	// 	fmt.Errorf(result.Error.Error())
	// }
	// userf := &User{}
	// result = db.Find(&userf)
	// fmt.Println(userf.Name)
}

// 初始化表结构，如果表不存在，自动创建，存在则不进行操作
func InitTable() {
	// 初始化用户账号表
	global.DB.AutoMigrate(
		&usersModel.User{},
		&constantModel.Constant{},
		&constantModel.CstDetail{},
		&constantModel.EventCst{},
		&gatewayModel.Gateway{},
	)
}

// 数据库配置
var databaseConf config.Database

// 初始化连接，返回连接
func gormSqlite() *gorm.DB {
	databaseConf = global.Conf.Database
	db, err := gorm.Open(sqlite.Open(databaseConf.FileName), &gorm.Config{})
	if err != nil {
		log.Println("Gorm sqlite start err: ", err)
		return nil
	}
	return db
}
