package respository

import "time"

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
