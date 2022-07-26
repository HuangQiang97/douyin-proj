package service

import (
	"douyin-proj/src/repository"
	"douyin-proj/src/types"
	"log"
	"sync"
)

// GetFeedVideos 根据截至时间获得视频
func GetFeedVideos(lastTime int64, isAuth bool, userId uint) (feedVideos []types.Video, nextTime int64, err error) {
	// 视频流
	videoList, err := repository.GetVideoTimeDesc(lastTime)
	if err != nil {
		log.Printf("根据截至时间获取视频流失败。lastTime:%d,err:%s\n", lastTime, err)
		return feedVideos, nextTime, err
	}
	// 返回的视频集合
	feedVideos = make([]types.Video, len(videoList), len(videoList))
	// 不存在满足条件的视频
	if len(videoList) == 0 {
		return feedVideos, lastTime, nil
	}
	// 下次拉取视频时的截至时间
	nextTime = int64(videoList[len(videoList)-1].CreatedAt)
	// go程
	var wg sync.WaitGroup
	wg.Add(len(videoList))
	// 填充视频信息
	for idx, video := range videoList {
		go func(idx int, video repository.Video) {
			user, _ := GetUserInfo(video.AuthorID, userId)
			feedVideo := types.Video{
				Id:            video.ID,
				PlayUrl:       video.PlayUrl,
				CoverUrl:      video.CoverUrl,
				FavoriteCount: uint64(getVideoFavoriteCount(video.ID)),
				CommentCount:  uint64(getVideoCommentCount(video.ID)),
				Title:         video.Title,
				Author:        *user,
				IsFavorite:    isFavorite(video.ID, userId),
			}
			feedVideos[idx] = feedVideo
			wg.Done()
		}(idx, video)
	}
	wg.Wait()
	return feedVideos, nextTime, nil
}
