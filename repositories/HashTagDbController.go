package repositories

import (
	"log"
	"strings"
	"trader-web-api/models"

	"gorm.io/gorm"
)

type HashTagRepository interface {
	Find() ([]string, error)
	Update(name string) error
	Create(tag []string) error
}
type HashTagRepoImpl struct {
	orm *gorm.DB
}

func newHashtagRepository(orm *gorm.DB) HashTagRepository {
	return &HashTagRepoImpl{
		orm: orm,
	}
}
func (r *HashTagRepoImpl) Find() ([]string, error) {
	var tag models.Hashtag
	result := r.orm.First(&tag)
	var temp = strings.Split(tag.Tags, SPLITKEY)
	return temp, result.Error
}
func (r *HashTagRepoImpl) Create(tags []string) error {
	var tagModel = models.Hashtag{
		Tags: strings.Join(tags, SPLITKEY),
	}
	result := r.orm.Create(&tagModel)
	return result.Error
}

func (r *HashTagRepoImpl) Update(name string) error {
	var tag models.Hashtag
	result := r.orm.First(&tag)
	if result.Error != nil {
		return result.Error
	}
	log.Println(tag.ID)
	tag.Tags = name
	result = r.orm.Save(tag)
	return result.Error
}
