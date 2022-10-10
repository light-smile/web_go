package test

import (
	"dnds_go/global"
	constantModel "dnds_go/src/models/constant"
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
	"testing"
)

type FindArgs struct {
	Keys []string
	Vals []interface{}
}

var DB *gorm.DB

func InitDb() {
	db, err := gorm.Open(sqlite.Open("../sql.db"), &gorm.Config{})
	if err != nil {
		log.Println("Gorm sqlite start err: ", err)
		//return nil
	}
	DB = db
}
func FindCom(f FindArgs) *gorm.DB {
	var conditionStr string
	for i := 0; i < len(f.Keys); i++ {
		if i == 0 {
			conditionStr += f.Keys[i] + "= ?"
			continue
		}
		conditionStr += "And" + f.Keys[i] + "= ?"
	}
	global.L.Lock()
	defer global.L.Unlock()
	return DB.Where(conditionStr, f.Vals...)
}
func TestFindCom(t *testing.T) {
	InitDb()
	find := FindArgs{
		Keys: []string{
			"id",
		},
		Vals: []interface {
		}{
			"1",
		},
	}
	res := FindCom(find)

	var user []constantModel.EventCst
	res.Find(&user)
	fmt.Println(user)
}
func FindComById(v int) *gorm.DB {
	return DB.Where("id = ?", v)
}
func TestFindComById(t *testing.T) {
	InitDb()
	var user []constantModel.EventCst
	FindComById(1).Find(&user)
	fmt.Println(user)
}
