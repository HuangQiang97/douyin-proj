package service

import (
	"github.com/HuangQiang97/douyin-proj/dal"
	"github.com/HuangQiang97/douyin-proj/entity"
	"github.com/HuangQiang97/douyin-proj/pkg/constant"
	"github.com/HuangQiang97/douyin-proj/pkg/errno"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string,error){
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), constant.PassWordCost)
	if err != nil{
		return "", err
	}
	return string(hashPassword), nil
}

func CreateUser(name string, password string) error {
	existedUser , err := dal.GetUserByName(name)
	if err != nil{
		return err
	}
	if len(existedUser) > 0{
		return errno.UserAlreadyExistErr
	}
	hashPw, err := hashPassword(password)
	if err != nil{
		return err
	}
	userModel := dal.User{
		UserName: name,
		Password: hashPw,
	}
	return dal.CreateUser(&userModel)
}

func CheckUser(name string, password string) (int64,error) {
	existedUser, err := dal.GetUserByName(name)
	if err != nil{
		return 0,err
	}
	if len(existedUser) == 0{
		return 0,errno.UserNotExistErr
	}
	if err := bcrypt.CompareHashAndPassword([]byte(existedUser[0].Password), []byte(password)); err != nil{
		return 0,err
	}
	return int64(existedUser[0].ID), nil
}

func GetUserByName(name string)(*entity.User, error) {
	existedUser, err := dal.GetUserByName(name)
	if err != nil{
		return nil, err
	}
	res := &entity.User{
		Id: int64(existedUser[0].ID),
		Name: existedUser[0].UserName,
		FollowerCount: 0,
		FollowCount: 0,
		IsFollow: false,
	}
	return res, nil
}


func GetUserById(id int64) (*entity.User,error){
	existedUser, err := dal.GetUserById([]int64{id})
	if err != nil {
		return nil, err
	}
	res := &entity.User{
		Id:            int64(existedUser[0].ID),
		Name:          existedUser[0].UserName,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}
	return res, nil
}

func GetUserAllInfo(id int64)(*entity.User, error) {
	user , err := GetUserById(id)
	if err != nil{
		return nil, err
	}
	followList, err := dal.GetAllFollow(id)
	if err != nil{
		return nil, err
	}
	followerList, err := dal.GetAllFollower(id)
	user.FollowCount = int64(len(followList))
	user.FollowerCount = int64(len(followerList))
	return user, nil
}