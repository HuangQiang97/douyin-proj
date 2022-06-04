package types

type FavoriteRequest struct {
	Token      string `form:"token" binding:"required"`
	VideoId    uint   `form:"video_id" binding:"required"`
	ActionType uint8  `form:"action_type" binding:"required,oneof=1 2"`
}

type FavoriteResponse struct {
	Response
}

type FavoriteListRequest struct {
	UserId uint   `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}

type FavoriteListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
