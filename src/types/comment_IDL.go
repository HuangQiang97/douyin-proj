package types

// CommentRequest 评论请求bean
type CommentRequest struct {
	Token       string `form:"token" binding:"required"`
	VideoId     uint   `form:"video_id" binding:"required"`
	ActionType  uint8  `form:"action_type" binding:"required,oneof=1 2"`
	CommentText string `form:"comment_text" binding:"omitempty"`
	CommentId   uint   `form:"comment_id" binding:"omitempty"`
}

// CommentResponse 评论响应bean
type CommentResponse struct {
	Response
	Comment Comment `json:"Comment"`
}

// CommentListRequest 获取视频评论列表请求
type CommentListRequest struct {
	Token   string `form:"token" binding:"required"`
	VideoId uint   `form:"video_id" binding:"required"`
}

// CommentListResponse 获取视频评论列表响应
type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}
