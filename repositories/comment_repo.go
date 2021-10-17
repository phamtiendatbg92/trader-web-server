package repositories

import (
	"trader-web-api/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentRepository interface {
	CreateCommentLv1(models.Commentlv1) uint
	CreateCommentLv2(models.Commentlv2) (uint, bool)
	GetAllComment(postID uint) ([]models.Commentlv1, error)
	DeleteCommentLv1(id uint, userId uint) bool
	DeleteCommentLv2(id uint, userId uint) bool
}

type commentRepoImpl struct {
	orm *gorm.DB
}

func newCommentRepository(orm *gorm.DB) CommentRepository {
	return &commentRepoImpl{
		orm: orm,
	}
}
func (r *commentRepoImpl) DeleteCommentLv1(id uint, userId uint) bool {
	var temp models.Commentlv1
	result := r.orm.Where("id=?", id).Find(&temp)
	zap.S().Info("result: ", temp)
	if result.Error == nil && result.RowsAffected != 0 {
		if temp.UserID != userId {
			zap.S().Info("User is different")
			return false
		} else {
			result = r.orm.Select("Commentlv2s").Delete(&temp)
			if result.Error != nil {
				zap.S().Info("ERROR not nil")
				return false
			} else {
				return true
			}
		}
	}
	zap.S().Info("Default return false")
	return false
}
func (r *commentRepoImpl) DeleteCommentLv2(id uint, userId uint) bool {
	var temp models.Commentlv2
	result := r.orm.Where("id=?", id).Find(&temp)
	zap.S().Info("result: ", temp)
	if result.Error == nil && result.RowsAffected != 0 {
		if temp.UserID != userId {
			zap.S().Info("User is different")
			return false
		} else {
			result = r.orm.Delete(&models.Commentlv2{}, id)
			if result.Error != nil {
				zap.S().Info("ERROR not nil")
				return false
			} else {
				return true
			}
		}
	}
	zap.S().Info("Default return false")
	return false
}
func (r *commentRepoImpl) CreateCommentLv1(comment models.Commentlv1) uint {
	result := r.orm.Create(&comment)
	if result.Error != nil {
		return 0
	}
	if result.RowsAffected == 0 {
		return 0
	} else {
		return comment.ID
	}
}
func (r *commentRepoImpl) CreateCommentLv2(comment models.Commentlv2) (uint, bool) {
	result := r.orm.Create(&comment)
	if result.Error != nil {
		return 0, false
	}
	if result.RowsAffected == 0 {
		return 0, false
	} else {
		return comment.ID, true
	}
}
func (r *commentRepoImpl) GetAllComment(postID uint) ([]models.Commentlv1, error) {
	var temp []models.Commentlv1
	zap.S().Info("get all comment with post_id: ", postID)
	result := r.orm.Where("post_id=?", postID).Preload(clause.Associations).Find(&temp)
	return temp, result.Error
}
