package service

import (
	"crypto/md5"
	"douyin-proj/src/config"
	"douyin-proj/src/repository"
	"douyin-proj/src/types"
	"encoding/hex"
	"errors"
	"log"
)

// encryptPassword 加密用户密码
func encryptPassword(password string) string {
	pd := []byte(password)
	salt := []byte(config.Salt)
	h := md5.New()
	h.Write(salt) // 先写盐值
	h.Write(pd)
	return hex.EncodeToString(h.Sum(nil))
}

// CreateUser 创建用户
func CreateUser(username, password string) (uint, error) {
	// 数据校验
	if _, err := repository.GetUserByName(username); err == nil {
		return 0, errors.New("用户已存在")
	}
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

// CheckUser 登录检查
func CheckUser(username, password string) (uint, error) {
	user, err := repository.GetUserByName(username)
	if err != nil {
		return uint(config.UserNotExisted), err
	}
	p := encryptPassword(password)
	if user.Password != p {
		return uint(config.WrongPassword), errors.New("password error")
	}
	return user.ID, nil
}

// GetUserInfo 获取用户信息
func GetUserInfo(targetId uint, currId uint) (responseUser *types.User, err error) {
	// 数据校验
	if !repository.ExistUser(targetId) {
		return nil, errors.New("用户不存在")
	}
	// 基本信息
	basicUser, err := repository.GetUserById(targetId)
	if err != nil {
		log.Printf("获取用户信息失败。uid:%d,err:%s\n", targetId, err)
		return nil, err
	}
	// 返回响应用户信息
	responseUser = &types.User{
		Id:            basicUser.ID,
		Name:          basicUser.UserName,
		FollowCount:   uint64(GetFollowingCount(targetId)),
		FollowerCount: uint64(GetFollowerCount(targetId)),
		IsFollow:      IsFollowing(currId, targetId),
	}
	return responseUser, nil
}
