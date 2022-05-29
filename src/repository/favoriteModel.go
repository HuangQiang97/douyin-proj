package repository

type Favorite struct {
	UserID  uint `gorm:"primarykey"`
	VideoID uint `gorm:"primarykey"`
}

func (f *Favorite) TableName() string {
	return "favorite"
}

func CreateFavorite(f *Favorite) error {
	return DB.Create(f).Error
}

func UndoFavorite(f *Favorite) error {
	return DB.Delete(f).Error
}
