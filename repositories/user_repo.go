package repositories

import (
	"trader-web-api/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(userID uint) (*models.Account, error)
	FindByName(name string) (*models.Account, error)
	FindByEmail(email string) (*models.Account, error)
	Create(account models.Account) bool
}

type userRepoImpl struct {
	orm *gorm.DB
}

func newUserRepository(orm *gorm.DB) UserRepository {
	return &userRepoImpl{
		orm: orm,
	}
}
func (r *userRepoImpl) Create(account models.Account) bool {
	result := r.orm.Create(&account)
	if result.Error != nil {
		return false
	} else {
		return true
	}
}
func (r *userRepoImpl) FindByID(userID uint) (*models.Account, error) {
	var res models.Account
	err := r.orm.Model(&models.Account{}).Where("id = ?", userID).First(&res).Error
	return &res, err
}

func (r *userRepoImpl) FindByName(name string) (*models.Account, error) {
	var res models.Account
	err := r.orm.Model(&models.Account{}).Where("name = ?", name).First(&res).Error
	return &res, err
}
func (r *userRepoImpl) FindByEmail(email string) (*models.Account, error) {
	var res models.Account
	err := r.orm.Model(&models.Account{}).Where("email = ?", email).First(&res).Error
	return &res, err
}
