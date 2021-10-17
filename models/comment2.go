package models

type Commentlv2 struct {
	ID       uint `gorm:"primaryKey"`
	UserID   uint
	Comment  string
	ParentId uint
}

func (Commentlv2) TableName() string {
	return "commentlv2"
}
