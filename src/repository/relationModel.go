package repository

type Relation struct {
	ID       uint `gorm:"primarykey"`
	UserID   uint
	FollowID uint
}

func (r *Relation) TableName() string {
	return "relation"
}

func CreateRelation(r *Relation) error {
	return DB.Create(r).Error
}
