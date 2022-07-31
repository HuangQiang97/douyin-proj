package types

// FavoriteRequest 点赞请求
type FavoriteRequest struct {
	Token      string `form:"token" binding:"required"`
	VideoId    uint   `form:"video_id" binding:"required"`
	ActionType uint8  `form:"action_type" binding:"required,oneof=1 2"`
}

// FavoriteResponse 点赞响应
type FavoriteResponse struct {
	Response
}

// FavoriteListRequest 获取用户点赞视频列表请求
type FavoriteListRequest struct {
	UserId uint   `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}

// FavoriteListResponse 获取用户点赞视频列表响应
type FavoriteListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
