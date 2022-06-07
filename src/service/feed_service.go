package service

import (
	"douyin-proj/src/global/util"
	"douyin-proj/src/repository"
	"douyin-proj/src/types"
	"log"
)

// GetFeedVideos 根据截至时间获得视频
func GetFeedVideos(lastTime int64, isAuth bool, uid uint) (feedVideos []types.Video, nextTime int64, err error) {
	// 视频流
	videoIds, err := repository.GetVideoIdsTimeDesc(lastTime)
	videoList := make([]repository.Video, 0, len(videoIds))
	for _, vid := range videoIds {
		if util.VideoExist(uint(vid)) {
			video, _ := util.GetVideo(uint(vid))
			videoList = append(videoList, *video)
		} else {
			video, _ := repository.GetVideoById(uint(vid))
			videoList = append(videoList, *video)
			if video != nil {
				util.AddVideo(video)
			}
		}
	}

	//videoList, err := repository.GetVideoTimeDesc(lastTime)

	if err != nil {
		log.Printf("根据截至时间获取视频流失败。lastTime:%d,err:%s\n", lastTime, err)
		return feedVideos, nextTime, err
	}

	// 返回的视频集合
	feedVideos = make([]types.Video, 0, len(videoList))
	// 不存在满足条件的视频
	if len(videoList) == 0 {
		return feedVideos, lastTime, nil
	}
	// 下次拉取视频时的截至时间
	nextTime = int64(videoList[len(videoList)-1].CreatedAt)
	for _, video := range videoList {
		// 填充视频作者信息
		authorId := video.AuthorID
		var user *repository.User
		if util.UserExist(authorId) {
			user, _ = util.GetUser(authorId)
		} else {
			user, _ = repository.GetUserById(authorId)
			if user != nil {
				util.AddUser(user)
			}
		}
		//user, _ := repository.GetUserById(authorId)

		feedUser := types.User{
			Id:            user.ID,
			Name:          user.UserName,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FansCount,
			IsFollow:      false,
		}
		// 填充视频信息
		feedVideo := types.Video{
			Id:            video.ID,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			Title:         video.Title,
		}
		// 如果用户已经登录，获取视频点赞和关注信息
		if isAuth {
			// 是否关注视频作者
			if repository.GetRelation(&repository.Relation{UserID: uid, FollowID: authorId}) {
				feedUser.IsFollow = true
			}
			// 是否点赞该视频
			if repository.IsFavorite(&repository.Favorite{UserID: uid, VideoID: video.ID}) {
				feedVideo.IsFavorite = true
			}
		}
		feedVideo.Author = feedUser
		feedVideos = append(feedVideos, feedVideo)
	}
	return feedVideos, nextTime, nil
}
