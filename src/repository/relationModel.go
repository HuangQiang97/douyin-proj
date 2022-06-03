package repository

import (
	"errors"
	"gorm.io/gorm"
)

type Relation struct {
	UserID   uint `gorm:"primarykey"`
	FollowID uint `gorm:"primarykey"`
}

func (r *Relation) TableName() string {
	return "relation"
}

func CreateRelation(r *Relation) error {
	return DB.Create(r).Error
}

func GetRelation(r *Relation) bool {
	count := int64(0)
	DB.Table("relation").Where("user_id=? and follow_id=?", r.UserID, r.FollowID).Count(&count)
	return count > 0
}

func CreateRelationWithCount(userId, followId uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&Relation{UserID: userId, FollowID: followId}).Error; err != nil {
			return err
		}

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
		if db.RowsAffected != 2 {
			return errors.New("user not exist")
		}
		return nil
	})
}

func DeleteRelationWithCount(userId, followId uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("user_id = ? AND follow_id = ?", userId, followId).Delete(&Relation{})
		if err := db.Error; err != nil {
			return err
		}
		if db.RowsAffected != 1 {
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
			return errors.New("user not exist")
		}
		return nil
	})
}

func GetFollow(userId uint) ([]User, error) {
	var users []User
	err := DB.Table("user").Joins("join relation as r on r.follow_id = user.id  and r.user_id = ? ", userId).Find(&users).Error
	return users, err
}

func GetFans(userId uint) ([]User, error) {
	var users []User
	err := DB.Table("user").Joins("join relation as r on r.user_id = user.id  and r.follow_id = ? ", userId).Find(&users).Error
	return users, err
}
