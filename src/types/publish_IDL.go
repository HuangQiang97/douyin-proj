package types

import "mime/multipart"

// PublishRequest 发布视频请求
type PublishRequest struct {
	Data  *multipart.FileHeader `form:"data" binding:"required"`
	Token string                `form:"token" binding:"required"`
	Title string                `form:"title" binding:"required"`
}

type PublishResponse Response

// VideoListRequest 获取用户发表视频列表请求
type VideoListRequest struct {
	UserId int64  `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}

// VideoListResponse 获取用户发表视频列表响应
type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
