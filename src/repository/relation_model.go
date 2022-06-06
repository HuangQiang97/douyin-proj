package repository

import (
	"errors"
	"gorm.io/gorm"
	"log"
)

type Relation struct {
	UserID   uint `gorm:"primarykey"`
	FollowID uint `gorm:"primarykey"`
}

func (r *Relation) TableName() string {
	return "relation"
}

//func CreateRelation(r *Relation) error {
//	return DB.Create(r).Error
//}

// GetRelation 是否存在点赞关系
func GetRelation(r *Relation) bool {
	count := int64(0)
	DB.Table("relation").Where("user_id=? and follow_id=?", r.UserID, r.FollowID).Count(&count)
	return count > 0
}

// CreateRelationWithCount 添加关注关系
func CreateRelationWithCount(userId, followId uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// uid+followId作为联合主键，防止重复关注
		if err := tx.Create(&Relation{UserID: userId, FollowID: followId}).Error; err != nil {
			log.Printf("添加关注失败。uid:%d,followId:%d,err:%s\n", userId, followId, err)
			return err
		}

		// 更新关注数
		db := tx.Table("user").Where("id IN ?", []uint{userId, followId})
		db.Updates(map[string]interface{}{
			"follow_count": gorm.Expr(`CASE id
											WHEN ? THEN follow_count + 1
										    WHEN ? THEN follow_count + 0
										    END`, userId, followId),
			"fans_count": gorm.Expr(`CASE id
											WHEN ? THEN fans_count + 1
										    WHEN ? THEN fans_count + 0
										    END`, followId, userId),
		})
		if db.Error != nil {
			return db.Error
		}
		// 如果被关注用户不存在或者自己关注自己 db.RowsAffected != 2 将无法满足，触发回滚
		if db.RowsAffected != 2 {
			log.Printf("添加关注失败，关注用户不存在或者自己关注自己。uid:%d,followId:%d\n", userId, followId)
			return errors.New("user not exist")
		}
		return nil
	})
}

// DeleteRelationWithCount 删除关注信息
func DeleteRelationWithCount(userId, followId uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("user_id = ? AND follow_id = ?", userId, followId).Delete(&Relation{})
		if err := db.Error; err != nil {
			return err
		}
		// 如果关注信息不存在，将回滚，不进行关注数和被关注数更新
		if db.RowsAffected != 1 {
			log.Printf("删除关注失败，关注关系不存在。uid:%d,followId:%d\n", userId, followId)
			return errors.New("relation is not existed")
		}

		db = tx.Table("user").Where("id IN ?", []uint{userId, followId})
		db.Updates(map[string]interface{}{
			"follow_count": gorm.Expr(`CASE id
										WHEN ? THEN follow_count - 1
									    WHEN ? THEN follow_count + 0
									    END`, userId, followId),
			"fans_count": gorm.Expr(`CASE id
										WHEN ? THEN fans_count - 1
									    WHEN ? THEN fans_count + 0
									    END`, followId, userId),
		})
		if db.Error != nil {
			return db.Error
		}
		if db.RowsAffected != 2 {
			log.Printf("删除关注失败。uid:%d,followId:%d\n", userId, followId)
			return errors.New("user not exist")
		}
		return nil
	})
}

// GetFollow 获得关注用户
func GetFollow(userId uint) ([]User, error) {
	var users []User
	err := DB.Table("user").Joins("join relation as r on r.follow_id = user.id  and r.user_id = ? ", userId).Find(&users).Error
	return users, err
}

func GetFollowIds(userId uint) ([]int, error) {
	var ids []int
	err := DB.Table("relation").Where("user_id=?", userId).Select("follow_id").Find(&ids).Error
	return ids, err
}

func GetFollowerIds(userId uint) ([]int, error) {
	var ids []int
	err := DB.Table("relation").Where("follow_id=?", userId).Select("user_id").Find(&ids).Error
	return ids, err
}

// GetFans 获得粉丝
func GetFans(userId uint) ([]User, error) {
	var users []User
	err := DB.Table("user").Joins("join relation as r on r.user_id = user.id  and r.follow_id = ? ", userId).Find(&users).Error
	return users, err
}
