package repository

import "douyin-proj/src/database"

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

func CreateVideo(video *Video) error {
	return database.MySQLDb.Create(video).Error
}

func GetVideoById(id uint) (*Video, error) {
	video := Video{}
	if err := DB.First(&video, id).Error; err != nil {
		return nil, err
	}
	return &video, nil
}

func GetVideosByIds(ids []uint) ([]*Video, error) {
	var videos []*Video
	if err := DB.Find(&videos, ids).Error; err != nil {
		return nil, err
	}
	return videos, nil
}
