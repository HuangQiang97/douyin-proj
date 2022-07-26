package controller

import (
	"douyin-proj/src/config"
	"douyin-proj/src/server/middleware"
	"douyin-proj/src/service"
	"douyin-proj/src/types"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

// CommentAction 评论请求操作
func CommentAction(c *gin.Context) {
	// 参数绑定
	var commentRequest = types.CommentRequest{}
	if err := c.ShouldBind(&commentRequest); err != nil {
		log.Println("反序列化评论操作请求失败")
		c.JSON(http.StatusOK, types.CommentResponse{
			Response: types.Response{StatusCode: config.ParamInvalid, StatusMsg: err.Error()},
		})
		return
	}
	// 登录
	userId, err := middleware.VerifyToken(commentRequest.Token)
	if err != nil {
		log.Printf("登录失败。err:%s\n", err)
		c.JSON(http.StatusOK, types.CommentResponse{
			Response: config.AuthFailedResp,
		})
		return
	}
	switch commentRequest.ActionType {
	case 1: // 添加评论
		if len(strings.TrimSpace(commentRequest.CommentText)) == 0 {
			log.Println("评论内容为空")
			c.JSON(http.StatusOK, types.CommentResponse{
				Response: config.ParamInvalidResp,
			})
			return
		}
		// 创建评论
		comment, err := service.CreateComment(userId, commentRequest.VideoId, commentRequest.CommentText)
		if err != nil {
			c.JSON(http.StatusOK, types.CommentResponse{
				Response: config.CommentAddFailedResp,
			})
			return
		}
		c.JSON(http.StatusOK, types.CommentResponse{
			Response: config.SuccessResp,
			Comment:  *comment,
		})
		log.Printf("创建评论成功。commentId:%d,uid:%d,videoId:%d\n", comment.Id, userId, commentRequest.VideoId)
	// 删除评论
	case 2:
		user, err := service.DeleteCommentById(userId, commentRequest.VideoId, commentRequest.CommentId)
		if err != nil {
			c.JSON(http.StatusOK, types.CommentResponse{
				Response: config.CommentDeleteFailedResp,
			})
			return
		}
		c.JSON(http.StatusOK, types.CommentResponse{
			Response: config.SuccessResp,
			Comment:  types.Comment{User: *user},
		})
		log.Printf("删除评论成功。commentId:%d,uid:%d,videoId:%d\n", commentRequest.CommentId, userId, commentRequest.VideoId)
	}
}

// CommentList 按时间倒叙排列返回视频评论
func CommentList(c *gin.Context) {
	// 参数绑定
	var commentListRequest = types.CommentListRequest{}
	if err := c.ShouldBind(&commentListRequest); err != nil {
		log.Println("反序列化获取评论列表操作请求失败")
		c.JSON(http.StatusOK, types.CommentListResponse{
			Response: types.Response{StatusCode: config.ParamInvalid, StatusMsg: err.Error()},
		})
		return
	}
	// 鉴权
	userId, err := middleware.VerifyToken(commentListRequest.Token)
	if err != nil {
		log.Printf("登录失败，err:%s\n", err)
		c.JSON(http.StatusOK, types.CommentResponse{
			Response: config.AuthFailedResp,
		})
		return
	}
	// 获取评论
	commentList, err := service.GetCommentByVideoId(commentListRequest.VideoId, userId)
	if err != nil {
		c.JSON(http.StatusOK, types.CommentListResponse{
			Response: types.Response{StatusCode: config.UnknownError, StatusMsg: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, types.CommentListResponse{
		Response:    config.SuccessResp,
		CommentList: commentList,
	})
	log.Printf("获取评论列表成功。uid:%d,videoId:%d\n", userId, commentListRequest.VideoId)
}
