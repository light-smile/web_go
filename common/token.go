package common

import (
	loginModel "dnds_go/src/models/login"
	"time"
)

func GetLoginToken(u loginModel.LoginReq) loginModel.AccessToken {
	return loginModel.AccessToken{
		AccessToken:  "adladjflkdflkadf",
		RefreshToken: "adfladfkl",
		ExpresIn:     time.Now().AddDate(1, 0, 0).Unix(),
	}
}
