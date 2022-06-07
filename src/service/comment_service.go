package service

import (
	"douyin-proj/src/repository"
	"douyin-proj/src/types"
	"errors"
	"log"
	"time"
)

// CreateComment 创建用户评论
func CreateComment(userId uint, videoId uint, content string) (*types.Comment, error) {
	var comment = repository.Comment{
		UserID:     userId,
		VideoID:    videoId,
		Content:    content,
		CreateDate: uint64(time.Now().Unix()),
	}

	// 用户合法性判断
	user, err := repository.GetUserById(comment.UserID)
	if err != nil {
		log.Printf("获取用户失败。uid:%d,err:%s\n", comment.UserID, err)
		return nil, err
	}

	// 添加评论并更新视频评论数
	err = repository.CreateCommentWithCount(&comment)
	if err != nil {
		return nil, err
	}

	// 填充响应用户信息
	var commentResp = &types.Comment{
		Id: comment.ID,
		User: types.User{
			Id:            user.ID,
			Name:          user.UserName,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FansCount,
			IsFollow:      false,
		},
		Content:    comment.Content,
		CreateDate: time.Unix(int64(comment.CreateDate), 0).Format("01-02"), //"2006-01-02 15:04:01"
	}
	return commentResp, nil
}

// DeleteCommentById 根据评论ID删除评论，并将视频评论数减一
func DeleteCommentById(userId uint, videoId uint, commentId uint) (*types.User, error) {
	// 用户合法性判断
	user, err := repository.GetUserById(userId)
	if err != nil {
		log.Printf("获取用户失败。uid:%d,err:%s\n", userId, err)
		return nil, err
	}

	// 删除评论
	err = repository.DeleteComment(userId, videoId, commentId)
	if err != nil {
		return nil, err
	}
	// 填充响应用户信息
	var userResp = &types.User{
		Id:            user.ID,
		Name:          user.UserName,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FansCount,
		IsFollow:      false,
	}
	return userResp, nil

}

// GetCommentByVideoId 时间倒叙获取评论
func GetCommentByVideoId(videoId uint, uId uint) ([]types.Comment, error) {
	// 视频合法性判断
	if !repository.ExistVideo(&videoId) {
		return nil, errors.New("视频不存在")
	}
	// 获取评论
	comments, err := repository.GetCommentIdsByVideoId(videoId)
	if err != nil {
		log.Printf("获取视频评论失败。videoId:%d,err=%s\n", videoId, err)
		return nil, err
	}

	// 填充响应用户信息
	var commentResp = make([]types.Comment, 0, len(comments))
	for _, c := range comments {
		u, err := repository.GetUserById(c.UserID)
		if err != nil {
			log.Printf("获取用户失败。uid:%d ,err=%s\n", uId, err)
			continue
		}
		// 用户与评论作者关注关系
		isFollow := repository.GetRelation(&repository.Relation{UserID: uId, FollowID: c.UserID})
		comment := types.Comment{
			Id: c.ID,
			User: types.User{
				Id:            u.ID,
				Name:          u.UserName,
				FollowCount:   u.FollowCount,
				FollowerCount: u.FansCount,
				IsFollow:      isFollow,
			},
			Content:    c.Content,
			CreateDate: time.Unix(int64(c.CreateDate), 0).Format("01-02"),
		}
		commentResp = append(commentResp, comment)
	}
	return commentResp, nil
}
