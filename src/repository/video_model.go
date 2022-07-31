package repository

// Video 数据模型
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

// CreateVideo 创建视频
func CreateVideo(video *Video) error {
	return DB.Create(video).Error
}

// GetVideoById 根据ID获得视频
func GetVideoById(id uint) (*Video, error) {
	video := Video{}
	if err := DB.First(&video, id).Error; err != nil {
		return nil, err
	}
	return &video, nil
}

// GetVideoIdsByAuthorId 获取用户发布的全部视频
func GetVideoIdsByAuthorId(authorId uint) ([]uint, error) {
	var ids []uint
	err := DB.Table("video").Where("author_id = ?", authorId).Select("id").Find(&ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// GetVideoTimeDesc 倒叙时间获得视频
func GetVideoTimeDesc(lastTime int64) ([]Video, error) {
	var videoList []Video
	err := DB.Table("video").Where("created_at < ? ", lastTime).Order("created_at DESC").Limit(30).Find(&videoList).Error
	return videoList, err
}

// GetVideoIdsTimeDesc 根据实践倒叙获得视频ID
func GetVideoIdsTimeDesc(lastTime int64) ([]int, error) {
	var ids []int
	err := DB.Table("video").Where("created_at < ? ", lastTime).Order("created_at DESC").Limit(30).Select("id").Find(&ids).Error
	return ids, err
}

// ExistVideo 判断视频是否存在
func ExistVideo(id uint) bool {
	count := int64(0)
	DB.Table("video").Where("id=? ", id).Count(&count)
	return count > 0
}

// GetAllVideoIds 获得全部视频ID
func GetAllVideoIds() ([]uint, error) {
	var ids []uint
	err := DB.Table("video").Select("id").Find(&ids).Error
	return ids, err
}
