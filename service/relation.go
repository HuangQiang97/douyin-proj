package service

import "github.com/HuangQiang97/douyin-proj/dal"

func CheckRelation(userID int64, toUserID int64) (bool, error){
	isFollow, err := dal.CheckFollowRelation(userID, toUserID)
	if err != nil{
		return false, err
	}
	return isFollow, nil
}

