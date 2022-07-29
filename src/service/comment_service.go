package service

import (
	"douyin-proj/src/repository"
	"douyin-proj/src/server/middleware"
	"douyin-proj/src/types"
	"errors"
	"log"
	"sync"
	"time"
)

// CreateComment 创建用户评论
func CreateComment(userId uint, videoId uint, content string) (*types.Comment, error) {
	// 数据校验
	if !repository.ExistVideo(videoId) {
		return nil, errors.New("视频不存在")
	}
	// 创建评论
	var comment = repository.Comment{
		UserID:     userId,
		VideoID:    videoId,
		Content:    content,
		CreateDate: uint64(time.Now().Unix()),
	}
	// 添加评论并更新视频评论数
	err := repository.CreateCommentWithCount(&comment)
	if err != nil {
		return nil, err
	}
	// 获得发表评论用户信息
	user, err := GetUserInfo(userId, userId)
	if err != nil {
		log.Printf("获取用户失败。uid:%d,err:%s\n", userId, err)
		return nil, err
	}
	// 填充响应用户信息
	var commentResp = &types.Comment{
		Id:         comment.ID,
		User:       *user,
		Content:    comment.Content,
		CreateDate: time.Unix(int64(comment.CreateDate), 0).Format("01-02"), //"2006-01-02 15:04:01"
	}
	// 存入redis
	go func() {
		if middleware.ExistVideoComment(videoId) {
			middleware.AddVideoComment(videoId, comment.ID)
		}
	}()
	return commentResp, nil
}

// DeleteCommentById 根据评论ID删除评论
func DeleteCommentById(userId uint, videoId uint, commentId uint) (*types.User, error) {
	// 数据校验
	comment, _ := repository.GetCommentById(commentId)
	if comment == nil || comment.UserID != userId || comment.VideoID != videoId {
		return nil, errors.New("userId+videoId+commentId不合法")
	}

	// 删除评论
	err := repository.DeleteComment(userId, videoId, commentId)
	if err != nil {
		return nil, err
	}
	// 删除缓存
	if middleware.ExistVideoComment(videoId) {
		middleware.DelVideoComment(videoId, commentId)
	}
	// 填充响应用户信息
	userResp, _ := GetUserInfo(userId, userId)
	return userResp, nil
}

// GetCommentByVideoId 时间倒叙获取视频评论
func GetCommentByVideoId(videoId uint, userId uint) ([]types.Comment, error) {
	// 合法性判断
	if !repository.ExistVideo(videoId) {
		return nil, errors.New("视频不存在")
	}
	// 获取评论Id
	var commentIds []uint
	if middleware.ExistVideoComment(videoId) {
		commentIds, _ = middleware.GetVideoComment(videoId)
	} else {
		commentIds, _ = repository.GetCommentIdsByVideoId(videoId)
	}
	// Go程
	wg := &sync.WaitGroup{}
	wg.Add(len(commentIds))
	// 填充响应用户信息
	var commentsResp = make([]types.Comment, len(commentIds), len(commentIds))
	for idx, commentId := range commentIds {
		go func(idx int, commentId uint) {
			// 用户与评论作者关注关系
			basicComment, _ := repository.GetCommentById(commentId)
			user, _ := GetUserInfo(basicComment.UserID, userId)
			comment := types.Comment{
				Id:         basicComment.ID,
				User:       *user,
				Content:    basicComment.Content,
				CreateDate: time.Unix(int64(basicComment.CreateDate), 0).Format("01-02"),
			}
			commentsResp[idx] = comment
			wg.Done()
		}(idx, commentId)
	}
	wg.Wait()
	// 存入缓存
	go func() {
		if !middleware.ExistVideoComment(videoId) {
			for _, commentId := range commentIds {
				middleware.AddVideoComment(videoId, commentId)
			}
		}
	}()
	return commentsResp, nil
}
