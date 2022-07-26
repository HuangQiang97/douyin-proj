package service

import (
	"douyin-proj/src/repository"
	"douyin-proj/src/server/middleware"
	"douyin-proj/src/types"
	"errors"
	"log"
	"sync"
)

// CreateRelation 添加关注
func CreateRelation(userId, followId uint) error {
	// 添加关注到数据库
	if err := repository.CreateRelationWithCount(userId, followId); err != nil {
		return err
	}
	// 更新redis当前用户关注关系缓存
	if middleware.ExistUserFollowing(userId) {
		if err := middleware.AddUserFollowing(userId, followId); err != nil {
			return err
		}
	}
	// 更新redis被关注用户被关注关系患处
	if middleware.ExistUserFollower(followId) {
		if err := middleware.AddUserFollower(followId, userId); err != nil {
			return err
		}
	}
	return nil
}

// DeleteRelation 取消关注
func DeleteRelation(userId, followId uint) error {
	// 删除数据库中关注信息
	if err := repository.DeleteRelationWithCount(userId, followId); err != nil {
		return err
	}
	// 更新redis当前用户关注关系缓存
	if middleware.ExistUserFollowing(userId) {
		if err := middleware.DelUserFollowing(userId, followId); err != nil {
			return err
		}
	}
	// 更新redis被关注用户被关注关系患处
	if middleware.ExistUserFollower(followId) {
		if err := middleware.DelUserFollower(followId, userId); err != nil {
			return err
		}
	}
	return nil
}

// GetFollowingList 获取用户关注列表
func GetFollowingList(targetId, currId uint) ([]types.User, error) {
	// 用户合法性判断
	if !repository.ExistUser(targetId) {
		log.Printf("目标用户不存在。currId:%d\n", targetId)
		return nil, errors.New("目标用户不存在")
	}
	// 获得用户关注用户Id
	var followingIds []uint
	if middleware.ExistUserFollowing(targetId) {
		followingIds, _ = middleware.GetUserFollowing(targetId)
	} else {
		followingIds, _ = repository.GetFollowingIds(targetId)
		for _, followingId := range followingIds {
			middleware.AddUserFollowing(targetId, followingId)
		}
	}
	// go程
	var wg sync.WaitGroup
	wg.Add(len(followingIds))
	// 填充用户信息
	var followList = make([]types.User, len(followingIds), len(followingIds))
	for idx, followingId := range followingIds {
		go func(idx int, followingId uint) {
			user, _ := GetUserInfo(followingId, currId)
			followList[idx] = *user
			wg.Done()
		}(idx, followingId)
	}
	wg.Wait()
	return followList, nil
}

// GetFollowerList 获得用户粉丝列表
func GetFollowerList(targetId, currId uint) ([]types.User, error) {
	// 用户合法性判断
	if !repository.ExistUser(targetId) {
		log.Printf("目标用户不存在。currId:%d\n", targetId)
		return nil, errors.New("目标用户不存在")
	}
	// 获得用户粉丝Id
	var followerIds []uint
	if middleware.ExistUserFollower(targetId) {
		followerIds, _ = middleware.GetUserFollower(targetId)
	} else {
		followerIds, _ = repository.GetFollowerIds(targetId)
		for _, followerId := range followerIds {
			middleware.AddUserFollower(targetId, followerId)
		}
	}
	// go程
	var wg sync.WaitGroup
	wg.Add(len(followerIds))
	// 填充用户信息
	var followerList = make([]types.User, len(followerIds), len(followerIds))
	for idx, followerId := range followerIds {
		go func(idx int, followerId uint) {
			user, _ := GetUserInfo(followerId, currId)
			followerList[idx] = *user
			wg.Done()
		}(idx, followerId)
	}
	wg.Wait()
	return followerList, nil
}

// GetFollowingCount 获得用户专注用户数
func GetFollowingCount(userId uint) int64 {

	// redis中存在则从缓存中取
	if middleware.ExistUserFollowing(userId) {
		cnt, _ := middleware.CountUserFollowing(userId)
		return int64(cnt)
	}
	// 不存在则从数据库中取，并放入缓存
	followingIds, _ := repository.GetFollowingIds(userId)
	go func() {
		for _, followingId := range followingIds {
			middleware.AddUserFollowing(userId, followingId)
		}
	}()
	return int64(len(followingIds))
}

// GetFollowerCount 获得用户粉丝数
func GetFollowerCount(userId uint) uint {
	// redis中存在则从缓存中取
	if middleware.ExistUserFollower(userId) {
		cnt, _ := middleware.CountUserFollower(userId)
		return cnt
	}
	// 不存在则从数据库中取，并放入缓存
	followerIds, _ := repository.GetFollowerIds(userId)
	go func() {
		for _, followerId := range followerIds {
			middleware.AddUserFollower(userId, followerId)
		}
	}()
	return uint(len(followerIds))
}

// IsFollowing 两个用户将是否存在关注关系
func IsFollowing(currId uint, targetId uint) bool {
	// redis中存在则从缓存中取
	if middleware.ExistUserFollowing(currId) {
		return middleware.ExistFollowRelation(currId, targetId)
	}
	// 不存在则从数据库中取
	return repository.GetRelation(&repository.Relation{
		UserID:   currId,
		FollowID: targetId,
	})
}
