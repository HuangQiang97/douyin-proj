package repository

import (
	"errors"
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
	return DB.Create(&comment).Error
}

func CreateCommentWithCount(comment *Comment) error {
	comment.CreateDate = uint64(time.Now().Unix())
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(comment).Error; err != nil {
			return err
		}
		//更新评论数
		db := tx.Table("video").Where("id = ?", comment.VideoID).Update("comment_count", gorm.Expr("comment_count + 1"))
		if db.Error != nil || db.RowsAffected != 1 {
			return errors.New("create comment failed")
		}
		return nil
	})
}

func DeleteComment(userId uint, videoId uint, commentId uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("user_id = ? AND video_id = ?", userId, videoId).Delete(&Comment{}, commentId)
		if db.Error != nil || db.RowsAffected != 1 {
			return errors.New("delete comment failed")
		}
		//更新评论数
		db = tx.Table("video").Where("id=?", videoId).Update("comment_count", gorm.Expr("comment_count - 1"))
		if db.Error != nil || db.RowsAffected != 1 {
			return errors.New("delete comment failed")
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

func GetCommentsByIds(ids []uint) ([]*Comment, error) {
	var comments []*Comment
	if err := DB.Find(&comments, ids).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func GetCommentIdsByVideoId(videoId uint) ([]Comment, error) {
	var comments []Comment
	db := DB.Session(&gorm.Session{}).Table("comment").Where("video_id = ?", videoId).Order("create_date DESC").Find(&comments)
	if db.Error != nil {
		return nil, db.Error
	}
	return comments, nil
}
