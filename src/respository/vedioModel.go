package respository

import "time"

type Video struct {
	ID            uint `gorm:"primarykey"`
	AuthorID      uint
	PlayUrl       string
	CoverUrl      string
	FavoriteCount int64
	CommentCount  int64
	Title         string
	CreatedAt     time.Time
}

func (v *Video) TableName() string {
	return "video"
}
