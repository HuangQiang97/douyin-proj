package repository

import (
	"errors"
	"gorm.io/gorm"
)

type Favorite struct {
	UserID  uint `gorm:"primarykey"`
	VideoID uint `gorm:"primarykey"`
}

func (f *Favorite) TableName() string {
	return "favorite"
}

func UndoFavorite(f *Favorite) error {
	// if there is no favorite record, return error
	if !IsFavorite(f) {
		return errors.New("no such favorite record")
	}
	return DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		if err := tx.Delete(f).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}

		var favoriteCount uint64
		// fetch favorite count from table `video`
		if err := tx.Model(&Video{}).Select("favorite_count").Where("id=?", f.VideoID).Find(&favoriteCount).Error; err != nil {
			return err
		}

		favoriteCount -= 1

		// update favorite count to table `video`
		if err := tx.Model(&Video{}).Where("id=?", f.VideoID).Update("favorite_count", favoriteCount).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
}

func CreateFavorite(f *Favorite) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		if err := tx.Create(f).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}

		var favoriteCount uint64
		// fetch favorite count from table `video`
		if err := tx.Model(&Video{}).Select("favorite_count").Where("id=?", f.VideoID).Find(&favoriteCount).Error; err != nil {
			return err
		}

		favoriteCount += 1

		// update favorite count to table `video`
		if err := tx.Model(&Video{}).Where("id=?", f.VideoID).Update("favorite_count", favoriteCount).Error; err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
}

func InsertFavorite(f *Favorite) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		if err := tx.Create(f).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		db := tx.Table("video").Where("id = ?", f.VideoID).UpdateColumn("favorite_count", gorm.Expr("favorite_count+1"))
		if err := db.Error; err != nil {
			return err
		}
		if db.RowsAffected != 1 {
			return errors.New("video is not existed")
		}
		// 返回 nil 提交事务
		return nil
	})
}

func DeleteFavorite(f *Favorite) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Delete(f)
		if err := db.Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		if db.RowsAffected != 1 {
			return errors.New("favorite is not existed")
		}
		tx.Table("video").Where("id = ?", f.VideoID).UpdateColumn("favorite_count", gorm.Expr("favorite_count-1"))
		if err := tx.Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
}

func IsFavorite(f *Favorite) bool {
	return DB.Find(f).RowsAffected == 1
}

func GetFavoriteVideoIdsByUserId(userId uint) ([]uint, error) {
	var videoIds []uint
	db := DB.Model(&Favorite{}).Select("video_id").Where("user_id=?", userId).Find(&videoIds)
	if db.Error != nil {
		return nil, db.Error
	}
	// if db.RowsAffected == 0 {
	// 	return nil, errors.New("record not exist")
	// }
	return videoIds, nil
}

func GetFavoriteVideoByUserId(userId uint) ([]Video, error) {
	var videos []Video
	err := DB.Table("video").Joins("join favorite on favorite.video_id = video.id  and favorite.user_id = ?", userId).Find(&videos).Error
	return videos, err
}
