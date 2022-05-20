package dal

import (
	"github.com/HuangQiang97/douyin-proj/pkg/constant"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName  string `json:"user_name"`
	Password   string `json:"password"`
}

func (u *User) TableName() string{
	return constant.UserTableName
}


func CreateUser(user *User) error {
	return DB.Create(user).Error
}

func GetUserById(ids []int64)([]*User, error){
	var res []*User
	if len(ids) == 0{
		return nil, nil
	}
	err := DB.Model(&User{}).Where("id in ?",ids).Find(&res).Error
	if err != nil{
		return nil, err
	}
	return res, nil
}

func GetUserByName(name string)([]*User, error){
	var res []*User
	err := DB.Model(&User{}).Where("user_name = ?", name).Find(&res).Error
	if err != nil{
		return nil, err
	}
	return res ,nil
}
