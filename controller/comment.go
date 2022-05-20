package controller

import (
	"github.com/HuangQiang97/douyin-proj/entity"
	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []entity.Comment `json:"comment_list,omitempty"`
}


// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {

}
