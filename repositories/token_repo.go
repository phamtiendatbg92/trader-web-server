package repositories

import (
	"trader-web-api/models"

	"gorm.io/gorm"
)

type TokenRepository interface {
	Save(m *models.Token) error
	Delete(m *models.Token) error
	FindByRefreshToken(refreshToken string) (*models.Token, error)
	DeleteByUserID(userID uint) bool
}

type tokenRepoImpl struct {
	orm *gorm.DB
}

func newTokenRepository(orm *gorm.DB) TokenRepository {
	return &tokenRepoImpl{
		orm: orm,
	}
}

func (r *tokenRepoImpl) DeleteByUserID(userID uint) bool {
	result := r.orm.Where("user_id=?", userID).Delete(&models.Token{})
	if result.Error != nil {
		return false
	} else {
		return true
	}
}
func (r *tokenRepoImpl) Save(m *models.Token) error {
	result := r.orm.Where("user_id=?", m.UserID).First(&models.Token{})
	if result.RowsAffected != 0 {
		// update
		return r.orm.Model(&models.Token{}).Where("user_id=?", m.UserID).Update("refresh_token", m.RefreshToken).Error
	} else {
		return r.orm.Create(&m).Error
	}
}

func (r *tokenRepoImpl) Delete(m *models.Token) error {
	return r.orm.Delete(m).Error
}

func (r *tokenRepoImpl) FindByRefreshToken(refreshToken string) (*models.Token, error) {
	var res models.Token
	err := r.orm.Model(&models.Account{}).Where("refresh_token = ?", refreshToken).First(&res).Error
	return &res, err
}
