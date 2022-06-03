package service

import (
	"douyin-proj/src/repository"
	"douyin-proj/src/types"
	"time"
)

func CreateComment(userId uint, videoId uint, content string) ([]types.Comment, error) {
	var comment = repository.Comment{
		UserID:  userId,
		VideoID: videoId,
		Content: content,
	}
	user, err := repository.GetUserById(comment.UserID)
	if err != nil {
		return nil, err
	}
	err = repository.CreateCommentWithCount(&comment)
	if err != nil {
		return nil, err
	}

	var commentResp = types.Comment{
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
	return []types.Comment{commentResp}, nil
}

func DeleteCommentById(userId uint, videoId uint, commentId uint) ([]types.User, error) {
	user, err := repository.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	err = repository.DeleteComment(userId, videoId, commentId)
	if err != nil {
		return nil, err
	}
	var userResp = types.User{
		Id:            user.ID,
		Name:          user.UserName,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FansCount,
		IsFollow:      false,
	}
	return []types.User{userResp}, nil

}

func GetCommentByVideoId(videoId uint, uId uint) ([]types.Comment, error) {
	comments, err := repository.GetCommentIdsByVideoId(videoId)
	if err != nil {
		return nil, err
	}

	var commentResp = make([]types.Comment, 0, len(comments))
	for _, c := range comments {
		u, err := repository.GetUserById(c.UserID)
		if err != nil {
			continue
		}
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
