package types

type PublishResponse Response

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
