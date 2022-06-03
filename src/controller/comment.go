package controller

import (
	"douyin-proj/src/global/ErrNo"
	"douyin-proj/src/global/util"
	"douyin-proj/src/repository"
	"douyin-proj/src/service"
	"douyin-proj/src/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// make/delete comment (only for signed user)
func CommentAction(c *gin.Context) {

	// check param
	var commentRequest = types.CommentRequest{}
	if err := c.ShouldBind(&commentRequest); err != nil {
		c.JSON(http.StatusOK, types.CommentResponse{
			Response: ErrNo.ParamInvalidResp,
		})
		return
	}

	if commentRequest.Token == "" {
		c.JSON(http.StatusOK, types.CommentResponse{
			Response: ErrNo.NotSignedInResp,
		})
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
		commentId, err := service.CreateComment(userId, commentRequest.VideoId, commentRequest.CommentText)
		if err != nil {
			c.JSON(http.StatusOK, types.CommentResponse{
				Response: ErrNo.CommentAddFailedResp,
			})
			return
		} else {

			com, err := repository.GetCommentById(commentId)
			if err != nil {
				c.JSON(http.StatusOK, types.CommentResponse{
					Response: ErrNo.CommentAddFailedResp,
				})
				return
			}

			user, err := repository.GetUserById(com.UserID)
			if err != nil {
				c.JSON(http.StatusOK, types.CommentResponse{
					Response: ErrNo.CommentAddFailedResp,
				})
				return
			}

			comment := types.Comment{
				Id: com.ID,
				User: types.User{
					Id:            user.ID,
					Name:          user.UserName,
					FollowCount:   user.FollowCount,
					FollowerCount: user.FansCount,
					IsFollow:      false,
				},
				Content:    com.Content,
				CreateDate: time.Unix(int64(com.CreateDate), 0).Format("2006-01-02 15:04:01"),
			}

			c.JSON(http.StatusOK, types.CommentResponse{
				Response: ErrNo.SuccessResp,
				Comment:  comment,
			})
			return

		}
	case 2: //delete comment
		if err := service.DeleteCommentById(commentRequest.CommentId); err != nil {
			c.JSON(http.StatusOK, types.CommentResponse{
				Response: ErrNo.CommentDeleteFailedResp,
			})
			return
		} else {
			c.JSON(http.StatusOK, types.CommentResponse{
				Response: ErrNo.SuccessResp,
			})
			return
		}

	}

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	var commentListRequest = types.CommentListRequest{}
	if err := c.ShouldBind(&commentListRequest); err != nil {
		c.JSON(http.StatusOK, types.CommentListResponse{
			Response: ErrNo.ParamInvalidResp,
		})
		return
	}

	commentIds, err := repository.GetCommentIdsByVideoId(commentListRequest.VideoId)
	fmt.Println(commentIds)
	if err != nil {
		c.JSON(http.StatusOK, types.CommentListResponse{
			Response: ErrNo.UnknownErrorResp,
		})
		return
	} else {

		// return if no comment list
		if len(commentIds) == 0 {
			c.JSON(http.StatusOK, types.CommentListResponse{
				Response:    ErrNo.SuccessResp,
				CommentList: []types.Comment{},
			})
			return
		}

		commentPtrList, err := repository.GetCommentsByIds(commentIds)
		if err != nil {
			c.JSON(http.StatusOK, types.CommentListResponse{
				Response: ErrNo.UnknownErrorResp,
			})
			return
		}

		commentList := make([]types.Comment, len(commentPtrList))
		for _, commentPtr := range commentPtrList {
			user, err := repository.GetUserById(commentPtr.UserID)
			if err != nil {
				c.JSON(http.StatusOK, types.CommentListResponse{
					Response: ErrNo.UnknownErrorResp,
				})
				return
			}

			comment := types.Comment{
				Id: commentPtr.ID,
				User: types.User{
					Id:            user.ID,
					Name:          user.UserName,
					FollowCount:   user.FollowCount,
					FollowerCount: user.FansCount,
					IsFollow:      false,
				},
				Content:    commentPtr.Content,
				CreateDate: time.Unix(int64(commentPtr.CreateDate), 0).Format("2006-01-02 15:04:01"),
			}

			commentList = append(commentList, comment)
		}
		c.JSON(http.StatusOK, types.CommentListResponse{
			Response:    ErrNo.SuccessResp,
			CommentList: commentList,
		})
		return
	}
}
