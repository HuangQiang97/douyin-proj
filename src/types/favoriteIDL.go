package types

type FavoriteRequest struct {
	UserId     uint   `json:"user_id,omitempty"  binding:"required"`
	Token      string `json:"token"  binding:"required"`
	VideoId    uint   `json:"video_id" binding:"required"`
	ActionType uint8  `json:"action_type" binding:"required,min=1,max=2"`
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
