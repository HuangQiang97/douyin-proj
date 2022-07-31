package repository

import (
	"errors"
	"gorm.io/gorm"
)

// User 数据模型
type User struct {
	ID          uint   `gorm:"primarykey"`
	UserName    string `json:"user_name" gorm:"unique"`
	Password    string `json:"password"`
	FollowCount uint64 `json:"follow_count" gorm:"default:0"`
	FansCount   uint64 `json:"fans_count" gorm:"default:0"`
}

type UserResp struct {
	User     `gorm:"embedded"`
	isFollow bool
}

func (u *User) TableName() string {
	return "user"
}

// CreateUser 创建用户
// username数据中存在unique约束，防止重复
func CreateUser(user *User) error {
	return DB.Create(user).Error
}

// GetUserById 根据ID获得用户
func GetUserById(id uint) (*User, error) {
	user := User{}
	if err := DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByName 根据用户名获得用户
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

// ExistUser 判断ID指向用户是否存在
func ExistUser(id uint) bool {
	count := int64(0)
	DB.Table("user").Where("id=? ", id).Count(&count)
	return count > 0
}

// GetAllUserIds 获得全部用户ID
func GetAllUserIds() ([]uint, error) {
	var ids []uint
	err := DB.Table("user").Select("id").Find(&ids).Error
	return ids, err
}
