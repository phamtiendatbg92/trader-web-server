package services

import (
	"net/http"
	"time"
	"trader-web-api/dtos"
	"trader-web-api/models"
	"trader-web-api/repositories"
	"trader-web-api/utils"

	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type AuthService interface {
	Login(username, rawPass string) *dtos.LoginResponse
	RefreshToken(refreshToken string) *dtos.LoginResponse
	Register(username, rawPass string) dtos.Meta
	Logout(userID uint) dtos.Meta
}

type authServiceImpl struct {
	userRepo  repositories.UserRepository
	tokenRepo repositories.TokenRepository
	jwtHelper utils.JWTHelper
}

func newAuthService(userRepo repositories.UserRepository,
	tokenRepo repositories.TokenRepository,
	jwtHelper utils.JWTHelper) AuthService {

	return &authServiceImpl{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtHelper: jwtHelper,
	}
}
func (s *authServiceImpl) Logout(userID uint) dtos.Meta {
	s.tokenRepo.DeleteByUserID(userID)
	return dtos.Meta{
		Code: http.StatusOK,
	}
}
func (s *authServiceImpl) Register(username, rawPass string) dtos.Meta {

	salt, err := utils.GenerateSaltPassword()
	zap.S().Info("Generate salt: ", salt)
	if err != nil {
		return dtos.Meta{
			Code:    http.StatusInternalServerError,
			Message: "Cannot Generate Salt",
		}
	}
	var account = models.Account{
		Email:    username,
		Password: utils.HashPassword(rawPass, salt),
		Salt:     salt,
	}
	zap.S().Info("Create account: ", account)
	result := s.userRepo.Create(account)
	if result {
		return dtos.Meta{
			Code: http.StatusOK,
		}
	} else {
		return dtos.Meta{
			Code:    http.StatusInternalServerError,
			Message: "Cannot create account",
		}
	}
}
func (s *authServiceImpl) Login(username, rawPass string) *dtos.LoginResponse {
	user, err := s.userRepo.FindByEmail(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			zap.S().Warn("User not found")
			return &dtos.LoginResponse{
				Meta: &dtos.Meta{
					Code:    http.StatusNotFound,
					Message: "user/pass not correct",
				},
			}
		}
		zap.S().Error("error when find user by name, detail: ", err)
		return &dtos.LoginResponse{Meta: dtos.InternalServerErrorMeta}
	}
	// compare password
	if utils.HashPassword(rawPass, user.Salt) != user.Password {
		zap.S().Warn("Password not match")
		return &dtos.LoginResponse{
			Meta: &dtos.Meta{
				Code:    http.StatusNotFound,
				Message: "user/pass not correct",
			},
		}
	}

	// generate token
	accessToken, err := s.genNewAccessToken(user)
	if err != nil {
		zap.S().Errorf("Generating token for user: %v, err: %v", user.Name, err)
		return &dtos.LoginResponse{Meta: dtos.InternalServerErrorMeta}
	}
	refreshToken := s.genNewRefreshToken(user)

	err = s.tokenRepo.Save(&models.Token{
		UserID:       user.ID,
		RefreshToken: refreshToken,
	})
	if err != nil {
		zap.S().Errorf("error when save refresh token, detail: ", err)
		return &dtos.LoginResponse{Meta: dtos.InternalServerErrorMeta}
	}

	return &dtos.LoginResponse{
		Meta: dtos.SuccessMeta,
		Data: &dtos.LoginResponseData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
}

func (s *authServiceImpl) RefreshToken(oldToken string) *dtos.LoginResponse {
	oldRfToken, err := s.tokenRepo.FindByRefreshToken(oldToken)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			zap.S().Warn("Refresh token not found")
			return &dtos.LoginResponse{
				Meta: &dtos.Meta{
					Code:    http.StatusNotFound,
					Message: "refresh token not correct",
				},
			}
		}
		zap.S().Error("error when find refresh token, detail: ", err)
		return &dtos.LoginResponse{Meta: dtos.InternalServerErrorMeta}
	}

	err = s.tokenRepo.Delete(oldRfToken)
	if err != nil {
		zap.S().Error("error when delete old refresh token, detail: ", err)
		return &dtos.LoginResponse{Meta: dtos.InternalServerErrorMeta}
	}

	user, err := s.userRepo.FindByID(oldRfToken.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			zap.S().Warn("User not found")
			return &dtos.LoginResponse{
				Meta: &dtos.Meta{
					Code:    http.StatusNotFound,
					Message: "user/pass not correct",
				},
			}
		}
		zap.S().Error("error when find user by name, detail: ", err)
		return &dtos.LoginResponse{Meta: dtos.InternalServerErrorMeta}
	}

	// generate token
	newAccessToken, err := s.genNewAccessToken(user)
	if err != nil {
		zap.S().Errorf("Generating token for user: %v, err: %v", user.Name, err)
		return &dtos.LoginResponse{Meta: dtos.InternalServerErrorMeta}
	}
	newRefreshToken := s.genNewRefreshToken(user)

	err = s.tokenRepo.Save(&models.Token{
		UserID:       user.ID,
		RefreshToken: newRefreshToken,
	})
	if err != nil {
		zap.S().Errorf("error when save refresh token, detail: ", err)
		return &dtos.LoginResponse{Meta: dtos.InternalServerErrorMeta}
	}

	return &dtos.LoginResponse{
		Meta: dtos.SuccessMeta,
		Data: &dtos.LoginResponseData{
			AccessToken:  newAccessToken,
			RefreshToken: newRefreshToken,
		},
	}
}

func (s *authServiceImpl) genNewAccessToken(user *models.Account) (string, error) {
	currentTime := time.Now()
	claims := dtos.AuthClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: currentTime.Add(10 * time.Hour).Unix(),
			IssuedAt:  currentTime.Unix(),
		},
		UserID:   user.ID,
		UserName: user.Name,
	}

	accessToken, err := s.jwtHelper.GenerateToken(&claims)
	if err != nil {
		zap.S().Errorf("Generating token for user: %v, err: %v", user.Name, err)
		return "", err
	}

	return accessToken, nil
}

func (s *authServiceImpl) genNewRefreshToken(user *models.Account) string {
	return utils.RandomRefreshToken()
}
