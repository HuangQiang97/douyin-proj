package service

import (
	"douyin-proj/src/repository"
	"douyin-proj/src/types"
)

func AddFavorite(userId uint, videoId uint) error {
	err := repository.InsertFavorite(&repository.Favorite{UserID: userId, VideoID: videoId})
	return err
}

func UndoFavorite(userId uint, videoId uint) error {
	err := repository.DeleteFavorite(&repository.Favorite{UserID: userId, VideoID: videoId})
	return err
}

func GetFavoriteVideoListByUserId(qId, uId uint) ([]types.Video, error) {
	videos, err := repository.GetFavoriteVideoByUserId(qId)
	if err != nil {
		return nil, err
	}
	var videoList = make([]types.Video, 0, len(videos))
	for _, v := range videos {
		video := types.Video{
			Id:            v.ID,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			Title:         v.Title,
		}
		isFavorite := repository.IsFavorite(&repository.Favorite{UserID: uId, VideoID: v.ID})
		video.IsFavorite = isFavorite
		u, err := repository.GetUserById(v.AuthorID)
		if err != nil {
			return nil, err
		}
		isFollow := repository.GetRelation(&repository.Relation{UserID: uId, FollowID: v.AuthorID})
		video.Author = types.User{
			Id:            u.ID,
			Name:          u.UserName,
			FollowCount:   u.FollowCount,
			FollowerCount: u.FansCount,
			IsFollow:      isFollow,
		}
		videoList = append(videoList, video)
	}
	return videoList, nil
}
