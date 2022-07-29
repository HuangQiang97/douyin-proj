package repository

import (
	"errors"
	"gorm.io/gorm"
)

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

func GetUserById(id uint) (*User, error) {
	user := User{}
	if err := DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUsersByIds(ids []uint) ([]User, error) {
	var users []User
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

func GetUserResponse(qid uint, uid uint) (user *User, isFollow bool) {
	subquery1 := DB.Table("user").Where("id = ?", qid).Select("*")
	subquery2 := DB.Table("relation").Where("user_id = ? AND follow_id = ?", uid, qid).Select("count(1) as is_follow")
	var userresp = UserResp{}
	DB.Table("(?) as u , (?) as r", subquery1, subquery2).Find(&userresp)
	return &userresp.User, userresp.isFollow
}

func ExistUser(id uint) bool {
	count := int64(0)
	DB.Table("user").Where("id=? ", id).Count(&count)
	return count > 0
}
