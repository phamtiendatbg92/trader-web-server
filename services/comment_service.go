package services

import (
	"net/http"
	"trader-web-api/dtos"
	"trader-web-api/models"
	"trader-web-api/repositories"

	"go.uber.org/zap"
)

type CommentService interface {
	GetAllComment(postId uint) ([]dtos.CommentItemRes, dtos.Meta)
	PushComment(userID uint, comment dtos.Comment1Req) uint
	PushReply(userID uint, comment dtos.Comment2Req) (uint, bool)
	DeleteComment(userID uint, comment dtos.DeleteCommentReq) bool
}

type commentServiceImpl struct {
	commentRepo repositories.CommentRepository
}

func newCommentService(commentRepo repositories.CommentRepository) CommentService {
	return &commentServiceImpl{
		commentRepo: commentRepo,
	}
}
func (s *commentServiceImpl) DeleteComment(userID uint, comment dtos.DeleteCommentReq) bool {
	if comment.Level == 1 {
		return s.commentRepo.DeleteCommentLv1(comment.ID, userID)
	} else {
		return s.commentRepo.DeleteCommentLv2(comment.ID, userID)
	}
}
func (s *commentServiceImpl) PushComment(userID uint, comment dtos.Comment1Req) uint {
	return s.commentRepo.CreateCommentLv1(models.Commentlv1{
		UserID:  userID, // TODO get user_id from access token
		Comment: comment.Comment,
		PostId:  comment.PostId,
	})
}
func (s *commentServiceImpl) PushReply(userID uint, comment dtos.Comment2Req) (uint, bool) {
	return s.commentRepo.CreateCommentLv2(models.Commentlv2{
		UserID:   userID, // TODO get user_id from access token
		Comment:  comment.Comment,
		ParentId: comment.ParentId,
	})
}

func (s *commentServiceImpl) GetAllComment(postId uint) ([]dtos.CommentItemRes, dtos.Meta) {
	zap.S().Info(" ============= ", postId)
	allComment, err := s.commentRepo.GetAllComment(postId)
	if err != nil {
		return nil, dtos.Meta{
			Code:    http.StatusInternalServerError,
			Message: "cannot get comment from DB",
		}
	}
	zap.S().Info("All comment lv1: ", allComment)
	var result []dtos.CommentItemRes
	for _, value := range allComment {
		var item dtos.CommentItemRes
		item.Id = value.ID
		item.CommentItem = dtos.Comment{
			UserId:  value.UserID,
			Comment: value.Comment,
		}

		for i := 0; i < len(value.Commentlv2s); i++ {
			var children = dtos.Comment{
				UserId:  value.Commentlv2s[i].UserID,
				Comment: value.Commentlv2s[i].Comment,
			}
			item.Replies = append(item.Replies, children)
		}
		result = append(result, item)
	}

	return result, dtos.Meta{Code: 200, Message: ""}
}
