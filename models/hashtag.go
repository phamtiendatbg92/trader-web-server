package models

type Hashtag struct {
	ID   uint   `gorm:"column:Id"`
	Tags string `gorm:"column:tags"`
}

func (Hashtag) TableName() string {
	return "hashtag"
}
