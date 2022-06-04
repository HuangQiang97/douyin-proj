package service

import (
	"douyin-proj/src/repository"
	"douyin-proj/src/types"
)

func CreateRelation(userId, followId uint) error {
	err := repository.CreateRelationWithCount(userId, followId)
	return err
}

func DeleteRelation(userId, followId uint) error {
	err := repository.DeleteRelationWithCount(userId, followId)
	return err
}

func GetFollowList(userId, uId uint) ([]types.User, error) {
	users, err := repository.GetFollow(userId)
	if err != nil {
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
	users, err := repository.GetFans(userId)
	if err != nil {
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
