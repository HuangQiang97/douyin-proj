package service

import (
	"crypto/md5"
	"douyin-proj/src/global/ErrNo"
	"douyin-proj/src/repository"
	"douyin-proj/src/types"
	"errors"
	"fmt"
)

func encryptPassword(password string) string {
	p := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", p)
}

func CreateUser(username, password string) (uint, error) {
	p := encryptPassword(password)
	var user = &repository.User{
		UserName: username,
		Password: p,
	}
	err := repository.CreateUser(user)
	if err != nil {
		return uint(0), err
	}
	return user.ID, nil
}

func CheckUser(username, password string) (uint, error) {
	user, err := repository.GetUserByName(username)
	if err != nil {
		return uint(ErrNo.UserNotExisted), err
	}
	p := encryptPassword(password)
	if user.Password != p {
		return uint(ErrNo.WrongPassword), errors.New("password error")
	}
	return user.ID, nil
}

func GetUserInfo(userId uint, id uint) (user *types.User, err error) {
	u, err := repository.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	isFollow := repository.GetRelation(&repository.Relation{UserID: id, FollowID: userId})
	user = &types.User{
		Id:            u.ID,
		Name:          u.UserName,
		FollowCount:   u.FollowCount,
		FollowerCount: u.FansCount,
		IsFollow:      isFollow,
	}
	return user, nil
}
