package service

import (
	"douyin-proj/src/global/util"
	"douyin-proj/src/repository"
	"douyin-proj/src/types"
	"errors"
	"log"
)

// AddFavorite 添加点赞
func AddFavorite(userId uint, videoId uint) error {

	err := repository.InsertFavorite(&repository.Favorite{UserID: userId, VideoID: videoId})
	if err == nil && util.VideoExist(videoId) {
		err = util.VideoFavoIncr(videoId)
	}
	return err
}

//UndoFavorite 取消点赞
func UndoFavorite(userId uint, videoId uint) error {
	err := repository.DeleteFavorite(&repository.Favorite{UserID: userId, VideoID: videoId})
	if err == nil && util.VideoExist(videoId) {
		err = util.VideoFavoDecr(videoId)
	}
	return err
}

// GetFavoriteVideoListByUserId 根据qid获得用户点赞过视频,uId为当前登录用户
func GetFavoriteVideoListByUserId(qId, uId uint) ([]types.Video, error) {

	// 用户合法性判断
	if !repository.ExistUser(qId) {
		log.Printf("目标用户不存在。uId:%d\n", qId)
		return nil, errors.New("目标用户不存在")
	}

	videoIds, err := repository.GetFavoriteVideoIdsByUserId(qId)
	if err != nil {
		log.Printf("获取用户点赞过视频id失败。uId:%d,err:%s\n", qId, err)
		return nil, err
	}
	videos := make([]repository.Video, 0, len(videoIds))
	for _, vid := range videoIds {
		var video *repository.Video
		if util.VideoExist(vid) {
			video, _ = util.GetVideo(vid)
		} else {
			video, _ = repository.GetVideoById(vid)
			if video != nil {
				util.AddVideo(video)
			}
		}
		videos = append(videos, *video)
	}

	//videos, err := repository.GetFavoriteVideoByUserId(qId)
	//if err != nil {
	//	log.Printf("获取用户点赞过视频失败。uId:%d,err:%s\n", qId, err)
	//	return nil, err
	//}

	// 填充响应视频信息
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

		// 视频作者信息
		var u *repository.User
		if util.UserExist(v.AuthorID) {
			u, err = util.GetUser(v.AuthorID)
		} else {
			u, err = repository.GetUserById(v.AuthorID)
			if u != nil {
				util.AddUser(u)
			}
		}
		//u, err := repository.GetUserById(v.AuthorID)

		if err != nil {
			log.Printf("获取用户信息失败。uId:%d,err:%s\n", v.AuthorID, err)
			return nil, err
		}
		// 当前用户与视频作者关注关系
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
