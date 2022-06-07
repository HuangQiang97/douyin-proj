package types

type CommentRequest struct {
	Token       string `form:"token" binding:"required"`
	VideoId     uint   `form:"video_id" binding:"required"`
	ActionType  uint8  `form:"action_type" binding:"required,oneof=1 2"`
	CommentText string `form:"comment_text" binding:"omitempty"`
	CommentId   uint   `form:"comment_id" binding:"omitempty"`
}

type CommentResponse struct {
	Response
	Comment Comment `json:"Comment"`
}

type CommentListRequest struct {
	Token   string `form:"token" binding:"required"`
	VideoId uint   `form:"video_id" binding:"required"`
}

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}
