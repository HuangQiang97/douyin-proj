package types

type RelationRequest struct {
	Token      string `form:"token" binding:"required"`
	ToUserId   int64  `form:"to_user_id" binding:"required"`
	ActionType int32  `form:"action_type" binding:"required,oneof=1 2"`
}

type RelationResponse struct {
	Response
}

type UserFollowListRequest struct {
	UserId int64  `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}

type UserFollowListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

type UserFansListRequest struct {
	UserId int64  `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}

type UserFansListResponse struct {
	Response
	UserList []User `json:"user_list"`
}
