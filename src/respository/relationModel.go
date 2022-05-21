package respository

type Relation struct {
	ID       uint `gorm:"primarykey"`
	UserID   uint
	FollowID uint
}

func (r *Relation) TableName() string {
	return "relation"
}
