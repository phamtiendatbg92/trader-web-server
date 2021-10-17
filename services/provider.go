package services

import (
	"trader-web-api/repositories"
	"trader-web-api/utils"
)

type ServiceProvider interface {
	GetAuthService() AuthService
	GetUserService() UserService
	GetTutService() TutorialService
	GetCommentService() CommentService
}

type serviceProviderImpl struct {
	authService     AuthService
	userService     UserService
	tutorialService TutorialService
	commentService  CommentService
}

func NewServiceProvider(repoProvider repositories.RepositoryProvider, jwtHelper utils.JWTHelper) ServiceProvider {
	authService := newAuthService(repoProvider.GetUserRepo(),
		repoProvider.GetTokenRepo(),
		jwtHelper)
	userService := newUserService(repoProvider.GetUserRepo())

	tutService := newTutorialService(repoProvider.GetTutRepo(), repoProvider.GetHashTagRepo())
	commentService := newCommentService(repoProvider.GetCommentRepo())
	return &serviceProviderImpl{
		authService:     authService,
		userService:     userService,
		tutorialService: tutService,
		commentService:  commentService,
	}
}
func (s serviceProviderImpl) GetCommentService() CommentService {
	return s.commentService
}
func (s serviceProviderImpl) GetTutService() TutorialService {
	return s.tutorialService
}
func (s serviceProviderImpl) GetAuthService() AuthService {
	return s.authService
}

func (s serviceProviderImpl) GetUserService() UserService {
	return s.userService
}
