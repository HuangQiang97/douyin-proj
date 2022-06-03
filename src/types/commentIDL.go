package types

type CommentRequest struct {
	Token       string `json:"token"    binding:"required"`
	VideoId     uint   `json:"video_id"    binding:"required"`
	ActionType  uint8  `json:"action_type"    binding:"required,min=1,max=2"`
	CommentText string `json:"comment_text"    binding:"optional"`
	CommentId   uint   `json:"comment_id"    binding:"optional"`
}

type CommentListRequest struct {
	Token   string `json:"token"    binding:"required"`
	VideoId uint   `json:"video_id"    binding:"required"`
}

type CommentResponse struct {
	Response
	Comment Comment `json:"Comment"`
}

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}
