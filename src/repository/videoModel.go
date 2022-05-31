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

type VideoResp struct {
	Video      `gorm:"embedded"`
	isFavorite bool `gorm:"column:is_favorite"`
}

func (v *Video) TableName() string {
	return "video"
}

func CreateVideo(video *Video) error {
	return DB.Create(video).Error
}

func GetVideoById(id uint) (*Video, error) {
	video := Video{}
	if err := DB.First(&video, id).Error; err != nil {
		return nil, err
	}
	return &video, nil
}

func GetVideosByIds(ids []uint) ([]Video, error) {
	var videos []Video
	if err := DB.Find(&videos, ids).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

func GetVideoByAuthorId(authorId uint) ([]Video, error) {
	var videos []Video
	err := DB.Where("author_id = ?", authorId).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func GetVideoByAuthorIdWithFavorite(authorId uint, id uint) []VideoResp {
	var videolist = []VideoResp{}
	subquery := DB.Table("favorite").Where("user_id = ? AND video_id = video.id", id).Select("count(1)")
	//DB.Table("video").Where("author_id = ?", authorId).Select("*,(?) as is_favorite", subquery).Find(&videolist)
	rows, _ := DB.Table("video").Where("author_id = ?", authorId).Select("*,(?) as is_favorite", subquery).Rows()
	for rows.Next() {
		v := VideoResp{}
		rows.Scan(&v.ID, &v.AuthorID, &v.PlayUrl, &v.CoverUrl, &v.FavoriteCount, &v.CommentCount, &v.Title, &v.CreatedAt, &v.isFavorite)
		videolist = append(videolist, v)
	}
	return videolist
}
