package repository

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID         uint `gorm:"primarykey"`
	UserID     uint
	VideoID    uint
	Content    string
	CreateDate uint64
}

func (c *Comment) TableName() string {
	return "comment"
}

func CreateComment(comment *Comment) error {
	comment.CreateDate = uint64(time.Now().Unix())
	return DB.Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Create(comment).Error; err != nil {
				return err
			}

			var commentCount uint64
			//获取评论数
			if err := tx.Model(&Video{}).Select("comment_count").Where("id=?", comment.VideoID).Find(&commentCount).Error; err != nil {
				return err
			}

			commentCount += 1

			//更新评论数
			if err := tx.Model(&Video{}).Where("id=?", comment.VideoID).Update("comment_count", commentCount).Error; err != nil {
				return err
			}

			return nil
		})
}

func GetCommentById(id uint) (*Comment, error) {
	comment := Comment{}
	if err := DB.First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func GetCommentIdsByVideoId(videoId uint) ([]uint, error) {
	var commentIds []uint
	db := DB.Session(&gorm.Session{}).Table("comment").Select("comment_id").Where("video_id=?", videoId).Find(&commentIds)
	if db.Error != nil {
		return nil, db.Error
	}
	return commentIds, nil
}

func GetCommentsByIds(ids []uint) ([]*Comment, error) {
	var comments []*Comment
	if err := DB.Find(&comments, ids).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func DeleteComment(comment *Comment) error {
	return DB.Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Delete(comment).Error; err != nil {
				return err
			}

			var commentCount uint64
			//获取评论数
			if err := tx.Model(&Video{}).Select("comment_count").Where("id=?", comment.VideoID).Find(&commentCount).Error; err != nil {
				return err
			}

			commentCount -= 1

			//更新评论数
			if err := tx.Model(&Video{}).Where("id=?", comment.VideoID).Update("comment_count", commentCount).Error; err != nil {
				return err
			}

			return nil
		})
}
