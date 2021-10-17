package models

type Token struct {
	UserID       uint   `gorm:"column:user_id"`
	RefreshToken string `gorm:"column:refresh_token"`
}

func (Token) TableName() string {
	return "token"
}
