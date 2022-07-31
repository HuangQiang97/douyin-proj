package types

// FeedResponse 视频流拉取请求
type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// FeedRequest 视频流拉取响应
type FeedRequest struct {
	Token    string `form:"token,,omitempty"`
	LastTime int64  `json:"last_time,omitempty"`
}
