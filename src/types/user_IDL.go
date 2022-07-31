package types

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	UserName string `form:"username" json:"username" binding:"required,max=32"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=32"`
}

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

// UserInfoRequest 用户信息获取请求
type UserInfoRequest struct {
	UserId uint   `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}

// UserResponse 用户信息获取响应
type UserResponse struct {
	Response
	User User `json:"user"`
}
