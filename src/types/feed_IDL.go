package types

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

type FeedRequest struct {
	Token    string `form:"token,,omitempty"`
	LastTime int64  `json:"last_time,omitempty"`
}
