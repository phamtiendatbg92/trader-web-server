package repositories

import "gorm.io/gorm"

type RepositoryProvider interface {
	GetUserRepo() UserRepository
	GetTokenRepo() TokenRepository
	GetTutRepo() TutorialRepository
	GetHashTagRepo() HashTagRepository
	GetCommentRepo() CommentRepository
}

type repoProviderImpl struct {
	userRepo    UserRepository
	tokenRepo   TokenRepository
	tutRepo     TutorialRepository
	hashtagRepo HashTagRepository
	commentRepo CommentRepository
}

func NewRepositoryProvider(db *gorm.DB) (RepositoryProvider, error) {
	userRepo := newUserRepository(db)
	tokenRepo := newTokenRepository(db)
	tutRepo := newTutorialRepository(db)
	hashtagRepo := newHashtagRepository(db)
	commentRepo := newCommentRepository(db)
	return &repoProviderImpl{
		userRepo:    userRepo,
		tokenRepo:   tokenRepo,
		tutRepo:     tutRepo,
		hashtagRepo: hashtagRepo,
		commentRepo: commentRepo,
	}, nil
}
func (r *repoProviderImpl) GetCommentRepo() CommentRepository {
	return r.commentRepo
}

func (r *repoProviderImpl) GetTutRepo() TutorialRepository {
	return r.tutRepo
}
func (r *repoProviderImpl) GetUserRepo() UserRepository {
	return r.userRepo
}

func (r *repoProviderImpl) GetTokenRepo() TokenRepository {
	return r.tokenRepo
}
func (r *repoProviderImpl) GetHashTagRepo() HashTagRepository {
	return r.hashtagRepo
}
