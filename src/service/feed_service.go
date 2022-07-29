package service

import (
	"douyin-proj/src/repository"
	"douyin-proj/src/types"
	"log"
	"math"
	"sort"
	"sync"
)

// RateIdPair 视频评分、ID对
type RateIdPair struct {
	rate float64
	id   uint
}

type RateIdPairs []RateIdPair

func (s RateIdPairs) Len() int {
	return len(s)
}

// Less 降序排列
func (s RateIdPairs) Less(i, j int) bool {
	return s[i].rate > s[j].rate
}

func (s RateIdPairs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// GetFeedVideosWithRecommendation 获得视频流
func GetFeedVideosWithRecommendation(lastTime int64, isAuth bool, userId uint) (feedVideos []types.Video, nextTime int64, err error) {
	var videoList []repository.Video
	// 未登录，无法进行推荐
	if !isAuth {
		videoList, err = repository.GetVideoTimeDesc(lastTime)
		if err != nil {
			log.Printf("根据截至时间获取视频流失败。lastTime:%d,err:%s\n", lastTime, err)
			return feedVideos, nextTime, err
		}
	} else {
		videoIds := cfRecommendation(userId)
		videoList = make([]repository.Video, len(videoIds), len(videoIds))
		for idx, videoId := range videoIds {
			video, _ := repository.GetVideoById(videoId)
			videoList[idx] = *video
		}
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

//recommendation 协同过滤推荐算法
func cfRecommendation(targetUserId uint) []uint {
	// 全部用户与视频
	allUserIds, _ := repository.GetAllUserIds()
	allVideoIds, _ := repository.GetAllVideoIds()
	// 用户ID以及视频ID可能不连续，需要手动映射:id->idx
	userIdMap := make(map[uint]int, len(allUserIds))
	videoIdMap := make(map[uint]int, len(allVideoIds))
	videoIdMapRev := make(map[int]uint, len(allVideoIds))
	for idx, userId := range allUserIds {
		userIdMap[userId] = idx
	}
	for idx, videoId := range allVideoIds {
		videoIdMap[videoId] = idx
		videoIdMapRev[idx] = videoId
	}
	// 构建评分矩阵
	rateMatrix := make([][]int, len(allUserIds), len(allUserIds))
	for _, userId := range allUserIds {
		videoIds, _ := repository.GetFavoriteVideoIdsByUserId(userId)
		rate := make([]int, len(allVideoIds), len(allVideoIds))
		for _, videoId := range videoIds {
			rate[videoIdMap[videoId]] = 1
		}
		rateMatrix[userIdMap[userId]] = rate
	}
	// 当前用户评分列表
	targetUserRate := rateMatrix[userIdMap[targetUserId]]
	// 计算当前用户与其它用户相似度
	simList := make([]float64, len(allUserIds), len(allUserIds))
	for _, userId := range allUserIds {
		currUserRate := rateMatrix[userIdMap[userId]]
		var dot float64
		squareSum0 := 1e-9
		squareSum1 := 1e-9
		// 计算目标用户与当前用户相似度
		for i := 0; i < len(allVideoIds); i++ {
			dot += float64(targetUserRate[i] * currUserRate[i])
			squareSum0 += float64(targetUserRate[i] * targetUserRate[i])
			squareSum1 += float64(currUserRate[i] * currUserRate[i])
		}
		norm := math.Sqrt(squareSum0) * math.Sqrt(squareSum1)
		simList[userIdMap[userId]] = dot / norm
	}
	// 计算目标用户对全部视频的可能评分
	videoRate := make([]float64, len(allVideoIds), len(allVideoIds))
	for _, userId := range allUserIds {
		// 目标用户与当前用户相似度
		sim := simList[userIdMap[userId]]
		// 当前用户评分列表
		rate := rateMatrix[userIdMap[userId]]
		for _, videoId := range allVideoIds {
			videoRate[videoIdMap[videoId]] += sim * float64(rate[videoIdMap[videoId]])
		}
	}
	// 根据估计的评分排序
	rateIdPairs := make(RateIdPairs, len(allVideoIds), len(allVideoIds))
	for idx, rate := range videoRate {
		rateIdPairs[idx] = RateIdPair{rate: rate, id: videoIdMapRev[idx]}
	}
	sort.Sort(rateIdPairs)
	videoIds, _ := repository.GetFavoriteVideoIdsByUserId(targetUserId)
	favoriteCount := len(videoIds)
	res := len(allVideoIds) - favoriteCount
	result := make([]uint, 0, res)
	// 去重
	set := make(map[uint]struct{}, favoriteCount)
	for _, id := range videoIds {
		set[id] = struct{}{}
	}
	for _, pair := range rateIdPairs {
		if _, ok := set[pair.id]; !ok {
			result = append(result, pair.id)
		}
	}
	return result
}
