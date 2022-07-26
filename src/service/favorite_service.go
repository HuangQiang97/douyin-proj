package service

import (
	"douyin-proj/src/repository"
	"douyin-proj/src/server/middleware"
	"douyin-proj/src/types"
	"errors"
	"log"
	"sync"
)

// AddFavorite 添加点赞
func AddFavorite(userId uint, videoId uint) error {
	// 添加点赞记录到数据库
	if err := repository.InsertFavorite(&repository.Favorite{UserID: userId, VideoID: videoId}); err != nil {
		return err
	}
	// 存入缓存
	if middleware.ExistUserFavorite(userId) {
		middleware.AddUserFavorite(userId, videoId)
	} else {
		videoIds, _ := repository.GetFavoriteVideoIdsByUserId(userId)
		for _, vId := range videoIds {
			middleware.AddUserFavorite(userId, vId)
		}
	}
	// 更新缓存中视频点赞数
	if middleware.ExistVideoFavoriteCount(videoId) {
		middleware.IncrVideoFavoriteCount(videoId)
	} else {
		middleware.InitVideoFavoriteCount(videoId, repository.GetVideoFavoriteCount(videoId))
	}
	return nil
}

//UndoFavorite 取消点赞
func UndoFavorite(userId uint, videoId uint) error {
	// 删除点赞数据库记录
	if err := repository.DeleteFavorite(&repository.Favorite{UserID: userId, VideoID: videoId}); err != nil {
		return err
	}
	// 删除缓存中用户点赞过视频数据
	if middleware.ExistUserFavorite(userId) {
		middleware.DelUserFavorite(userId, videoId)
	} else {
		videoIds, _ := repository.GetFavoriteVideoIdsByUserId(userId)
		for _, vId := range videoIds {
			middleware.AddUserFavorite(userId, vId)
		}
	}
	// 更新缓存中视频点赞数
	if middleware.ExistVideoFavoriteCount(videoId) {
		middleware.DecrVideoFavoriteCount(videoId)
	} else {
		middleware.InitVideoFavoriteCount(videoId, repository.GetVideoFavoriteCount(videoId))
	}
	return nil
}

// GetFavoriteVideoListByUserId 根据targetId获得用户点赞过视频,currId为当前登录用户
func GetFavoriteVideoListByUserId(targetId, currId uint) ([]types.Video, error) {
	// 用户合法性判断
	if !repository.ExistUser(targetId) {
		log.Printf("目标用户不存在。currId:%d\n", targetId)
		return nil, errors.New("目标用户不存在")
	}
	// 点赞过视频Id
	var videoIds []uint
	if middleware.ExistUserFavorite(targetId) {
		videoIds, _ = middleware.GetUserFavorite(targetId)
	} else {
		videoIds, _ = repository.GetFavoriteVideoIdsByUserId(targetId)
		for _, vId := range videoIds {
			middleware.AddUserFavorite(targetId, vId)
		}
	}
	// Go程
	var wg sync.WaitGroup
	wg.Add(len(videoIds))
	var videoList = make([]types.Video, len(videoIds), len(videoIds))
	// 填充视频信息
	for idx, videoId := range videoIds {
		go func(idx int, videoId uint) {
			basicVideo, _ := repository.GetVideoById(videoId)
			user, _ := GetUserInfo(basicVideo.AuthorID, currId)
			video := types.Video{
				Id:            basicVideo.ID,
				PlayUrl:       basicVideo.PlayUrl,
				CoverUrl:      basicVideo.CoverUrl,
				FavoriteCount: uint64(getVideoFavoriteCount(videoId)),
				CommentCount:  uint64(getVideoCommentCount(videoId)),
				Title:         basicVideo.Title,
				Author:        *user,
				IsFavorite:    isFavorite(videoId, currId),
			}
			videoList[idx] = video
			wg.Done()
		}(idx, videoId)
	}
	wg.Wait()
	return videoList, nil
}

// getVideoFavoriteCount 获得视频点赞数
func getVideoFavoriteCount(videoId uint) uint {
	// 尝试从缓存中获取
	if middleware.ExistVideoFavoriteCount(videoId) {
		cnt, _ := middleware.GetVideoFavoriteCount(videoId)
		return cnt
	} else {
		// 从数据库中获取
		cnt := repository.GetVideoFavoriteCount(videoId)
		middleware.InitVideoFavoriteCount(videoId, cnt)
		return uint(cnt)
	}
}

// getVideoCommentCount 获得视频评论数
func getVideoCommentCount(videoId uint) uint {
	// 尝试从缓存中获取
	if middleware.ExistVideoComment(videoId) {
		cnt, _ := middleware.CountVideoComment(videoId)
		return cnt
	} else {
		// 从数据库中获取
		commentIds, _ := repository.GetCommentIdsByVideoId(videoId)
		for _, commentId := range commentIds {
			middleware.AddVideoComment(videoId, commentId)
		}
		return uint(len(commentIds))
	}
}

// isFavorite 视频与用户间是否存在点赞关系
func isFavorite(videoId uint, userId uint) bool {
	// 尝试从缓存中获取
	if !middleware.ExistUserFavorite(userId) {
		videoIds, _ := repository.GetFavoriteVideoIdsByUserId(userId)
		for _, videoId := range videoIds {
			middleware.AddUserFavorite(userId, videoId)
		}
	}
	return middleware.ExistFavoriteRelation(userId, videoId)
}
