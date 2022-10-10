package testService

import (
	"dnds_go/global"
	"dnds_go/logger"
	"dnds_go/provider"
	testModle "dnds_go/src/models/test"
	"log"
	"testing"
)

func TestCreate(t *testing.T) {
	provider.InitConfig()
	provider.InitDatabase()
	// 初始化日志
	if err := logger.InitLogger(global.Conf.Logger); err != nil {
		log.Fatalln("logger init error: ", err)
	}
	cases := []struct {
		user *testModle.User
		want int
	}{
		{
			user: &testModle.User{
				Name: "test1",
				Age:  000,
			},

			want: 1,
		}, {
			user: &testModle.User{
				Name: "test2",
				Age:  000,
			},

			want: 1,
		}, {
			user: &testModle.User{
				Name: "test3",
				Age:  000,
			},

			want: 1,
		},
	}

	for _, cc := range cases {
		if err := cc.user.Create(); err != nil {
			t.Errorf("faild create user  %v: %v", cc.user, err)
		}
		var ress *[]testModle.User
		if err := cc.user.FindByName(ress); err != nil {
			t.Errorf("faild find user  %v: %v", cc.user, err)
		}

	}
}
