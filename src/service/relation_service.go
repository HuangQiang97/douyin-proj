package service

import (
	"douyin-proj/src/global/util"
	"douyin-proj/src/repository"
	"douyin-proj/src/types"
	"errors"
	"log"
)

func CreateRelation(userId, followId uint) error {
	err := repository.CreateRelationWithCount(userId, followId)
	if err == nil {
		if util.UserExist(userId) {
			util.UserFollowIncr(userId)
		}
		if util.UserExist(followId) {
			util.UserFollowerIncr(followId)
		}
	}
	return err
}

func DeleteRelation(userId, followId uint) error {
	err := repository.DeleteRelationWithCount(userId, followId)
	if err == nil {
		if util.UserExist(userId) {
			util.UserFollowDecr(userId)
		}
		if util.UserExist(followId) {
			util.UserFollowerDecr(followId)
		}
	}
	return err
}

func GetFollowList(userId, uId uint) ([]types.User, error) {

	// 用户合法性判断
	if !repository.ExistUser(userId) {
		log.Printf("目标用户不存在。uId:%d\n", userId)
		return nil, errors.New("目标用户不存在")
	}

	followIds, err := repository.GetFollowIds(userId)
	users := make([]repository.User, 0, len(followIds))
	for _, fId := range followIds {
		if util.UserExist(uint(fId)) {
			u, _ := util.GetUser(uint(fId))
			users = append(users, *u)
		} else {
			u, _ := repository.GetUserById(uint(fId))
			users = append(users, *u)
			if u != nil {
				util.AddUser(u)
			}
		}
	}
	//users, err := repository.GetFollow(userId)

	if err != nil {
		log.Printf("获取用户关注用户失败。uId:%d,err:%s\n", userId, err)
		return nil, err
	}
	var followList = make([]types.User, 0, len(users))
	for _, u := range users {
		isFollow := repository.GetRelation(&repository.Relation{UserID: uId, FollowID: u.ID})
		follow := types.User{
			Id:            u.ID,
			Name:          u.UserName,
			FollowCount:   u.FollowCount,
			FollowerCount: u.FansCount,
			IsFollow:      isFollow,
		}
		followList = append(followList, follow)
	}
	return followList, nil
}

func GetFansList(userId, uId uint) ([]types.User, error) {

	// 用户合法性判断
	if !repository.ExistUser(userId) {
		log.Printf("目标用户不存在。uId:%d\n", userId)
		return nil, errors.New("目标用户不存在")
	}

	followerIds, err := repository.GetFollowerIds(userId)
	users := make([]repository.User, 0, len(followerIds))
	for _, fId := range followerIds {
		if util.UserExist(uint(fId)) {
			u, _ := util.GetUser(uint(fId))
			users = append(users, *u)
		} else {
			u, _ := repository.GetUserById(uint(fId))
			users = append(users, *u)
			if u != nil {
				util.AddUser(u)
			}
		}
	}

	//users, err := repository.GetFans(userId)

	if err != nil {
		log.Printf("获取用户粉丝失败。uId:%d,err:%s\n", userId, err)
		return nil, err
	}
	var followList = make([]types.User, 0, len(users))
	for _, u := range users {
		isFollow := repository.GetRelation(&repository.Relation{UserID: uId, FollowID: u.ID})
		follow := types.User{
			Id:            u.ID,
			Name:          u.UserName,
			FollowCount:   u.FollowCount,
			FollowerCount: u.FansCount,
			IsFollow:      isFollow,
		}
		followList = append(followList, follow)
	}
	return followList, nil
}
