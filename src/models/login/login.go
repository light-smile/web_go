package loginModel

// 登录请求
type LoginReq struct {
	UserLogin    string `json:"userLogin" binding:"required"`
	UserPassWord string `json:"userPassword" binding:"required"`
}

// 登录响应
type LoginRes struct {
	AccessToken AccessToken `json:"accessToken"`
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpresIn     int64  `json:"expires_in"`
}
