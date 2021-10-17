package models

type Account struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Email    string
	Password string
	Salt     string
}

func (Account) TableName() string {
	return "account"
}
