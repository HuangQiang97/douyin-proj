package types

type UserLoginRequest struct {
	UserName string `form:"username" json:"username" binding:"required,max=30"`
	Password string `form:"password" json:"password" binding:"required,min=5,max=30"`
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserInfoRequest struct {
	UserId uint   `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}
