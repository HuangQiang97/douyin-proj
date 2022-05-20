package service

import (
	"crypto/md5"
	"douyin-proj/src/respository"
	"errors"
	"fmt"
)

func encryptPassword(password string) string {
	p := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", p)
}

func CreateUser(username, password string) (uint, error) {
	p := encryptPassword(password)
	var user = &respository.User{
		UserName: username,
		Password: p,
	}
	err := respository.CreateUser(user)
	if err != nil {
		return uint(0), err
	}
	return user.ID, nil
}

func CheckUser(username, password string) (uint, error) {
	user, err := respository.GetUserByName(username)
	if err != nil {
		return 0, err
	}
	p := encryptPassword(password)
	if user.Password != p {
		return 0, errors.New("password error")
	}
	return user.ID, nil
}
