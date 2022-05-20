package types

type CommentResponse Response

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}
