package service

import (
	"crypto/md5"
	"douyin-proj/src/config"
	"douyin-proj/src/global/ErrNo"
	"douyin-proj/src/global/util"
	"douyin-proj/src/repository"
	"douyin-proj/src/types"
	"encoding/hex"
	"errors"
	"log"
)

func encryptPassword(password string) string {
	pd := []byte(password)
	salt := []byte(config.Salt)
	h := md5.New()
	h.Write(salt) // 先写盐值
	h.Write(pd)
	return hex.EncodeToString(h.Sum(nil))
	//p := md5.Sum([]byte(password))
	//return fmt.Sprintf("%x", p)
}

func CreateUser(username, password string) (uint, error) {
	p := encryptPassword(password)
	var user = &repository.User{
		UserName: username,
		Password: p,
	}
	err := repository.CreateUser(user)
	if err != nil {
		log.Printf("创建用户失败。username:%s,err:%s\n", username, err)
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
	var u *repository.User
	if util.UserExist(userId) {
		u, err = util.GetUser(userId)
	} else {
		u, err = repository.GetUserById(userId)
		if u != nil {
			util.AddUser(u)
		}
	}
	//u, err := repository.GetUserById(userId)
	if err != nil {
		log.Printf("获取用户信息失败。uid:%d,err:%s\n", userId, err)
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
