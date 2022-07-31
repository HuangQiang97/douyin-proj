package types

// Response 响应公共部分
type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

func (r *Response) WithStatusMsg(str string) {
	r.StatusMsg = str
}

// User 用于响应的用户信息
type User struct {
	Id            uint   `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   uint64 `json:"follow_count"` // 删掉omitempty，否则如果是0，json不会序列化
	FollowerCount uint64 `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// Video 用于响应视频信息
type Video struct {
	Id            uint   `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount uint64 `json:"favorite_count"`
	CommentCount  uint64 `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title,omitempty"`
}

// Comment 用于响应的评论信息
type Comment struct {
	Id         uint   `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}
