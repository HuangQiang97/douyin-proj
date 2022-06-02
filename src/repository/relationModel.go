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

func GetRelation(r *Relation) bool {
	count := int64(0)
	DB.Table("relation").Where("user_id=? and follow_id=?", r.UserID, r.FollowID).Count(&count)
	return count > 0
}
