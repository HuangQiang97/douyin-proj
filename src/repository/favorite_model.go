package repository

import (
	"errors"
	"gorm.io/gorm"
	"log"
)

// Favorite 数据模型
type Favorite struct {
	UserID  uint `gorm:"primarykey"`
	VideoID uint `gorm:"primarykey"`
}

func (f *Favorite) TableName() string {
	return "favorite"
}

// InsertFavorite 添加点赞记录，并更新视频点赞数
func InsertFavorite(f *Favorite) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		// uid+video作为联合主键，防止重复点赞
		if err := tx.Create(f).Error; err != nil {
			// 返回任何错误都会回滚事务
			log.Printf("添加点赞记录失败。uid:%d,videoId:%d,err=%s\n", f.UserID, f.VideoID, err)
			return err
		}
		// 更新点赞数
		db := tx.Table("video").Where("id = ?", f.VideoID).UpdateColumn("favorite_count", gorm.Expr("favorite_count+1"))
		if err := db.Error; err != nil {
			log.Printf("更新视频点赞数失败。 uid:%d,videoId:%d,err=%s\n", f.UserID, f.VideoID, err)
			return err
		}
		// 如果db.RowsAffected != 1表示视频不存在，回滚之前插入的点赞记录
		if db.RowsAffected != 1 {
			log.Printf("添加的视频不存在。 videoId=%d\n", f.VideoID)
			return errors.New("video is not existed")
		}
		// 返回 nil 提交事务
		return nil
	})
}

// DeleteFavorite 删除点赞记录，并更新视频点赞数
func DeleteFavorite(f *Favorite) error {

	return DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Delete(f)
		if err := db.Error; err != nil {
			// 返回任何错误都会回滚事务
			log.Printf("删除点赞记录失败。uid:%d,videoId:%d,err=%s\n", f.UserID, f.VideoID, err)
			return err
		}
		// 如果db.RowsAffected != 1表示点赞记录不存在，后续点赞数减一不必执行
		if db.RowsAffected != 1 {
			log.Printf("点赞记录不存在。uid:%d,videoId:%d\n", f.UserID, f.VideoID)
			return errors.New("favorite is not existed")
		}
		tx.Table("video").Where("id = ?", f.VideoID).UpdateColumn("favorite_count", gorm.Expr("favorite_count-1"))
		if err := tx.Error; err != nil {
			log.Printf("视频点赞减一记录失败。videoId:%d,err:%s\n", f.VideoID, err)
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
}

// GetFavoriteVideoIdsByUserId 根据用户id获得点赞过视频id
func GetFavoriteVideoIdsByUserId(userId uint) ([]uint, error) {
	var videoIds []uint
	db := DB.Model(&Favorite{}).Select("video_id").Where("user_id=?", userId).Find(&videoIds)
	if db.Error != nil {
		return nil, db.Error
	}
	return videoIds, nil
}

// GetFavoriteVideoByUserId 根据uid获得用户点赞过视频
func GetFavoriteVideoByUserId(userId uint) ([]Video, error) {
	var videos []Video
	err := DB.Table("video").Joins("join favorite on favorite.video_id = video.id  and favorite.user_id = ?", userId).Find(&videos).Error
	return videos, err
}

// GetVideoFavoriteCount 根据videoId获得点赞频数
func GetVideoFavoriteCount(videoId uint) int64 {
	count := int64(0)
	DB.Table("favorite").Where(" video_id=?", videoId).Count(&count)
	return count
}

// ExistFavorite 判断是否存在点赞关系
func ExistFavorite(userId, videoId uint) bool {
	count := int64(0)
	DB.Table("favorite").Where("video_id=? and user_id=?", videoId, userId).Count(&count)
	return count > 0
}
