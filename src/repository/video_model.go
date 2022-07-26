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
	err := DB.Table("video").Where("author_id = ?", authorId).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func GetVideoIdsByAuthorId(authorId uint) ([]uint, error) {
	var ids []uint
	err := DB.Table("video").Where("author_id = ?", authorId).Select("id").Find(&ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func GetVideoByAuthorIdWithFavorite(authorId uint, id uint) []VideoResp {
	var videolist []VideoResp
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

// GetVideoTimeDesc 倒叙时间获得视频
func GetVideoTimeDesc(lastTime int64) ([]Video, error) {
	var videoList []Video
	err := DB.Table("video").Where("created_at < ? ", lastTime).Order("created_at DESC").Limit(30).Find(&videoList).Error
	return videoList, err
}

func GetVideoIdsTimeDesc(lastTime int64) ([]int, error) {
	var ids []int
	err := DB.Table("video").Where("created_at < ? ", lastTime).Order("created_at DESC").Limit(30).Select("id").Find(&ids).Error
	return ids, err
}

func ExistVideo(id *uint) bool {
	count := int64(0)
	DB.Table("video").Where("id=? ", id).Count(&count)
	return count > 0
}
