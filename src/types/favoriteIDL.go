package types

type FavoriteResponse Response

type FavoriteListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}