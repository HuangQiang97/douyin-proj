package dal

import (
	"errors"
	"github.com/HuangQiang97/douyin-proj/pkg/constant"
	"gorm.io/gorm"
)

type Relation struct {
	gorm.Model
	UserID    int64 `json:"user_id"`
	ToUserID  int64 `json:"to_user_id"`
}

func (r *Relation) TableNames() string{
	return constant.RelationTableName
}

func InsertRelation(r *Relation) error{
	if r.UserID <= 0 || r.ToUserID <= 0 {
		return errors.New("param is invalid")
	}
	isFollow, err := CheckFollowRelation(r.UserID,r.ToUserID)
	if err != nil{
		return err
	}
	if isFollow{
		return nil
	}
	return DB.Create(r).Error
}

func DeleteRelations(r *Relation) error{
	if r.UserID <= 0 || r.ToUserID <= 0 {
		return errors.New("param is invalid")
	}
	return DB.Delete(r).Error
}
// 返回userID关注的用户的ID列表
func GetAllFollow(UserID int64) ([]int64, error){
	res := make([]int64,0)
	if err := DB.Model(&Relation{}).Where("user_id = ?",UserID).Select("to_user_id").Find(&res).Error; err != nil{
		return nil, err
	}
	return res,nil
}

// 返回userID的粉丝ID列表
func GetAllFollower(UserID int64)([]int64, error){
	res := make([]int64,0)
	if err := DB.Model(&Relation{}).Where("to_user_id = ?",UserID).Select("user_id").Find(&res).Error; err != nil{
		return nil, err
	}
	return res, nil
}

func CheckFollowRelation(userID int64, toUserID int64) (bool, error){
	res := make([]*Relation,0)
	err := DB.Model(&Relation{}).Where("to_user_id = ? and user_id = ?" ,toUserID, userID).Find(&res).Error
	if err != nil{
		return false, err
	}
	if len(res) == 0{
		return false ,nil
	}
	return true, nil
}