package util

import (
	"douyin-proj/src/repository"
	"fmt"
	"testing"
)

func TestRedis(t *testing.T) {
	video := &repository.Video{
		ID:            1,
		AuthorID:      2,
		PlayUrl:       "abc",
		CoverUrl:      "def",
		FavoriteCount: 3,
		CommentCount:  4,
		Title:         "ghi",
		CreatedAt:     5,
	}

	AddVideo(video)

	VideoFavoIncr(video.ID)
	VideoFavoIncr(video.ID)
	VideoFavoDecr(video.ID)

	VideoCommIncr(video.ID)
	VideoCommIncr(video.ID)
	VideoCommIncr(video.ID)
	VideoCommIncr(video.ID)
	VideoCommDecr(video.ID)
	VideoCommDecr(video.ID)

	v, _ := GetVideo(video.ID)
	fmt.Printf("%+v\n", v)

	user := &repository.User{
		ID:          1,
		UserName:    "abc",
		Password:    "def",
		FollowCount: 2,
		FansCount:   3,
	}
	AddUser(user)
	UserFollowIncr(user.ID)
	UserFollowIncr(user.ID)
	UserFollowDecr(user.ID)

	UserFollowerIncr(user.ID)
	UserFollowerIncr(user.ID)
	UserFollowerIncr(user.ID)
	UserFollowerIncr(user.ID)
	UserFollowerDecr(user.ID)
	UserFollowerDecr(user.ID)

	u, _ := GetUser(user.ID)
	fmt.Printf("%+v\n", u)

	AddComments(1, &[]repository.Comment{
		{
			ID:         1,
			UserID:     2,
			VideoID:    3,
			Content:    "4",
			CreateDate: 5,
		},
		{
			ID:         2,
			UserID:     2,
			VideoID:    3,
			Content:    "4",
			CreateDate: 5,
		},
	})
	AddComment(1, &repository.Comment{
		ID:         3,
		UserID:     2,
		VideoID:    3,
		Content:    "4",
		CreateDate: 5,
	})
	DeleteComment(1, 2)
	c, _ := GetComments(1)
	fmt.Printf("%+v\n", c)
}
