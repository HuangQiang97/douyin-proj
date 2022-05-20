package respository

import (
	"douyin-proj/src/database"
	"errors"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName    string `json:"user_name" gorm:"unique"`
	Password    string `json:"password"`
	FollowCount uint   `json:"follow_count" gorm:"default:0"`
	FansCount   uint   `json:"fans_count" gorm:"default:0"`
}

func (u *User) TableName() string {
	return "user"
}

func CreateUser(user *User) error {
	return database.MySQLDb.Create(user).Error
}

func GetUserById(id uint) (*User, error) {
	user := User{}
	if err := DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUsersByIds(ids []uint) ([]*User, error) {
	var users []*User
	if err := DB.Find(&users, ids).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByName(username string) (*User, error) {
	var user = User{}
	db := DB.Session(&gorm.Session{}).Where("user_name=?", username).Find(&user)
	if db.Error != nil {
		return nil, db.Error
	}
	if db.RowsAffected == 0 {
		return nil, errors.New("user not exist")
	}
	return &user, nil
}
