package controller

import (
	"douyin-proj/src/global/ErrNo"
	"douyin-proj/src/global/util"
	"douyin-proj/src/service"
	"douyin-proj/src/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

// make/delete comment (only for signed user)
func CommentAction(c *gin.Context) {
	// check param
	var commentRequest = types.CommentRequest{}
	if err := c.ShouldBind(&commentRequest); err != nil {
		c.JSON(http.StatusOK, types.CommentResponse{
			Response: types.Response{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()},
		})
		return
	}

	userId, err := util.VerifyToken(commentRequest.Token)
	if err != nil {
		c.JSON(http.StatusOK, types.CommentResponse{
			Response: ErrNo.AuthFailedResp,
		})
		return
	}

	switch commentRequest.ActionType {
	case 1: //make comment
		if commentRequest.CommentText == "" {
			c.JSON(http.StatusOK, types.CommentResponse{
				Response: ErrNo.ParamInvalidResp,
			})
			return
		}
		comment, err := service.CreateComment(userId, commentRequest.VideoId, commentRequest.CommentText)
		if err != nil {
			c.JSON(http.StatusOK, types.CommentResponse{
				Response: ErrNo.CommentAddFailedResp,
			})
			return
		}

		c.JSON(http.StatusOK, types.CommentResponse{
			Response: ErrNo.SuccessResp,
			Comment:  comment[0],
		})

	case 2: //delete comment
		if commentRequest.CommentId == 0 {
			c.JSON(http.StatusOK, types.CommentResponse{
				Response: ErrNo.ParamInvalidResp,
			})
			return
		}
		user, err := service.DeleteCommentById(userId, commentRequest.VideoId, commentRequest.CommentId)
		if err != nil {
			c.JSON(http.StatusOK, types.CommentResponse{
				Response: ErrNo.CommentDeleteFailedResp,
			})
			return
		}
		c.JSON(http.StatusOK, types.CommentResponse{
			Response: ErrNo.SuccessResp,
			Comment:  types.Comment{User: user[0]},
		})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	var commentListRequest = types.CommentListRequest{}
	if err := c.ShouldBind(&commentListRequest); err != nil {
		c.JSON(http.StatusOK, types.CommentListResponse{
			Response: types.Response{StatusCode: ErrNo.ParamInvalid, StatusMsg: err.Error()},
		})
		return
	}

	userId, err := util.VerifyToken(commentListRequest.Token)
	if err != nil {
		c.JSON(http.StatusOK, types.CommentResponse{
			Response: ErrNo.AuthFailedResp,
		})
		return
	}

	commentList, err := service.GetCommentByVideoId(commentListRequest.VideoId, userId)
	if err != nil {
		c.JSON(http.StatusOK, types.CommentListResponse{
			Response: types.Response{StatusCode: ErrNo.UnknownError, StatusMsg: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, types.CommentListResponse{
		Response:    ErrNo.SuccessResp,
		CommentList: commentList,
	})
}
