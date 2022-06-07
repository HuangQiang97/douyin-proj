package repository

import (
	"errors"
	"gorm.io/gorm"
	"log"
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

//func CreateComment(comment *Comment) error {
//	comment.CreateDate = uint64(time.Now().Unix())
//	return DB.Create(&comment).Error
//}

// CreateCommentWithCount 添加评论并更新视频评论数
func CreateCommentWithCount(comment *Comment) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(comment).Error; err != nil {
			log.Printf("创建评论失败。err:%s\n", err)
			return err
		}
		//更新评论数
		db := tx.Table("video").Where("id = ?", comment.VideoID).Update("comment_count", gorm.Expr("comment_count + 1"))
		// 如果video_id不合法，将触发异常，事务回滚，撤销插入的评论
		if db.Error != nil || db.RowsAffected != 1 {
			log.Printf("更新视频点赞数失败。err:%s\n", db.Error)
			return errors.New("create comment failed")
		}
		return nil
	})
}

// DeleteComment 删除评论
func DeleteComment(userId uint, videoId uint, commentId uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("user_id = ? AND video_id = ?", userId, videoId).Delete(&Comment{}, commentId)
		if db.Error != nil {
			log.Printf("删除评论失败。userId=%d,commentId:%d,err:%s\n", userId, commentId, db.Error)
			return errors.New("delete comment failed")
		}
		// 如果userId+videoId+commentId不合法，将触发异常，后续的点赞数-1不会执行
		if db.RowsAffected != 1 {
			log.Printf("该评论不存在。userId=%d,commentId:%d,err:%s\n", userId, commentId, db.Error)
			return errors.New("delete comment failed")
		}

		//更新评论数
		db = tx.Table("video").Where("id=?", videoId).Update("comment_count", gorm.Expr("comment_count - 1"))
		if db.Error != nil || db.RowsAffected != 1 {
			log.Printf("点赞数减一失败。videoId:%d,err:%s\n", videoId, db.Error)
			return errors.New("delete comment failed")
		}
		return nil
	})
}

//func GetCommentById(id uint) (*Comment, error) {
//	comment := Comment{}
//	if err := DB.First(&comment, id).Error; err != nil {
//		return nil, err
//	}
//	return &comment, nil
//}

//func GetCommentsByIds(ids []uint) ([]*Comment, error) {
//	var comments []*Comment
//	if err := DB.Find(&comments, ids).Error; err != nil {
//		return nil, err
//	}
//	return comments, nil
//}

// GetCommentIdsByVideoId 获取评论
func GetCommentsByVideoId(videoId uint) ([]Comment, error) {
	var comments []Comment
	db := DB.Session(&gorm.Session{}).Table("comment").Where("video_id = ?", videoId).Order("create_date DESC").Find(&comments)
	if db.Error != nil {
		return nil, db.Error
	}
	return comments, nil
}

//func ExistComment(commentId *uint, uid *uint, videoId *uint) bool {
//	count := int64(0)
//	DB.Table("comment").Where("id=? and user_id=? and video_id=?", commentId, uid, videoId).Count(&count)
//	return count > 0
//}
