package types

// RelationRequest 点赞相关操作请求
type RelationRequest struct {
	Token      string `form:"token" binding:"required"`
	ToUserId   int64  `form:"to_user_id" binding:"required"`
	ActionType int32  `form:"action_type" binding:"required,oneof=1 2"`
}

// RelationResponse 点赞操作请求响应
type RelationResponse struct {
	Response
}

// UserFollowListRequest 获得关注用户列表请求
type UserFollowListRequest struct {
	UserId int64  `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}

// UserFollowListResponse 获得关注用户列表响应
type UserFollowListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// UserFansListRequest 获得粉丝列表请求
type UserFansListRequest struct {
	UserId int64  `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}

// UserFansListResponse 获得粉丝列表响应
type UserFansListResponse struct {
	Response
	UserList []User `json:"user_list"`
}
