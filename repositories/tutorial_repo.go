package repositories

import (
	"log"
	"trader-web-api/models"

	"gorm.io/gorm"
)

type TutorialRepository interface {
	FindAllTutorial() ([]models.Tutorial, error)
	CreateNewPost(tut models.Tutorial) error
	FindTutByUrl(url string) (models.Tutorial, error)
	UpdateTutorial(tut models.Tutorial) error
	DeleteTutorial(id uint) error
}
type tutorialRepoImpl struct {
	orm *gorm.DB
}

func newTutorialRepository(orm *gorm.DB) TutorialRepository {
	return &tutorialRepoImpl{
		orm: orm,
	}
}

var SPLITKEY = "SPLITKEY"

func (r *tutorialRepoImpl) FindAllTutorial() ([]models.Tutorial, error) {
	var res []models.Tutorial
	// result := r.orm.First(&res)
	result := r.orm.Find(&res)
	return res, result.Error
}

func (r *tutorialRepoImpl) CreateNewPost(tut models.Tutorial) error {
	result := r.orm.Create(&tut)
	return result.Error
}

func (r *tutorialRepoImpl) FindTutByUrl(url string) (models.Tutorial, error) {
	var temp models.Tutorial
	log.Print(url)
	result := r.orm.Where("url = ?", url).Find(&temp)
	log.Print(result.RowsAffected)
	return temp, result.Error
}

func (r *tutorialRepoImpl) UpdateTutorial(tut models.Tutorial) error {
	result := r.orm.Save(tut)
	return result.Error
}

func (r *tutorialRepoImpl) DeleteTutorial(id uint) error {
	result := r.orm.Delete(&models.Tutorial{}, id)
	return result.Error
}
