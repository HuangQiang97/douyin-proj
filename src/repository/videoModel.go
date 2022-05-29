package repository

type Video struct {
	ID            uint `gorm:"primarykey"`
	AuthorID      uint
	PlayUrl       string
	CoverUrl      string
	FavoriteCount uint64
	CommentCount  uint64
	Title         string
	CreatedAt     uint64
}

func (v *Video) TableName() string {
	return "video"
}
