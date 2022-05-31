package service

import (
	"douyin-proj/src/repository"
	"douyin-proj/src/types"
)

func GetVideoList(authorId uint) (videoList []types.Video, err error) {
	user, err := repository.GetUserById(authorId)
	if err != nil {
		return nil, err
	}
	videos, err := repository.GetVideoByAuthorId(authorId)
	if err != nil {
		return nil, err
	}

	var author = types.User{
		Id:            user.ID,
		Name:          user.UserName,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FansCount,
		IsFollow:      false,
	}
	videoList = make([]types.Video, 0, len(videos))
	for _, v := range videos {
		videoList = append(videoList, types.Video{
			Id:            v.ID,
			Author:        author,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    false,
			Title:         v.Title,
		})
	}
	return videoList, nil
}
