package models

type Commentlv1 struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	Comment     string
	PostId      uint
	Commentlv2s []Commentlv2 `gorm:"foreignKey:parent_id"`
}

func (Commentlv1) TableName() string {
	return "commentlv1"
}
