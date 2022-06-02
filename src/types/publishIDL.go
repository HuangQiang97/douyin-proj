package types

import "mime/multipart"

type PublishRequest struct {
	Data  *multipart.FileHeader `form:"data" binding:"required"`
	Token string                `form:"token" binding:"required"`
	Title string                `form:"title" binding:"required"`
}

type PublishResponse Response

type VideoListRequest struct {
	UserId int64  `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
