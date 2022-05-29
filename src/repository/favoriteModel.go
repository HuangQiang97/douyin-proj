package repository

import (
	"gorm.io/gorm"
)

type Favorite struct {
	UserID  uint `gorm:"primarykey"`
	VideoID uint `gorm:"primarykey"`
}

func (f *Favorite) TableName() string {
	return "favorite"
}

func CreateFavorite(f *Favorite) error {
	return DB.Create(f).Error
}

func UndoFavorite(f *Favorite) error {
	return DB.Delete(f).Error
}

func IsFavorite(f *Favorite) bool {
	return DB.Find(f).RowsAffected == 1
}

func GetFavoriteVideoIdsByUserId(userId uint) ([]uint, error) {
	var videoIds []uint
	db := DB.Session(&gorm.Session{}).Table("favorite").Select("video_id").Where("user_id=?", userId).Find(&videoIds)
	if db.Error != nil {
		return nil, db.Error
	}
	// if db.RowsAffected == 0 {
	// 	return nil, errors.New("record not exist")
	// }
	return videoIds, nil
}
