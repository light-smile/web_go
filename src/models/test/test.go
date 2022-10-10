package testModle

import (
	"dnds_go/global"
	"sync"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

/*
db = global.DB
	db.AutoMigrate(&user)
	db.Create(&user)
	var userres []testModle.User
	db.Find(&userres)
*/
var m sync.Mutex

func (u *User) Create() error {

	err := global.DB.Create(u).Error
	if err != nil {
		return err
	}
	return nil
}
func (u *User) FindByName(result *[]User) error {
	// fmt.Println(u.Name)
	if u.Name == "" {
		return nil
	}
	err := global.DB.Where("name = ?", u.Name).Find(&result).Error
	if err != nil {
		return err
	}
	// fmt.Println(result)
	return nil
}
