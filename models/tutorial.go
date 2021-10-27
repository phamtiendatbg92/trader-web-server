package models

type Tutorial struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	Content     string
	Url         string
	Tags        string // split by ','
	AuthorID    uint
	Description string
}

func (Tutorial) TableName() string {
	return "tutorial"
}
