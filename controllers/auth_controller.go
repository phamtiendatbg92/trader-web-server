package controllers

import (
	"net/http"
	"trader-web-api/dtos"
	"trader-web-api/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthController struct {
	BaseController
}

func NewAuthController(sp services.ServiceProvider) *AuthController {
	c := &AuthController{}
	c.serviceProvider = sp
	return c
}
func (c *AuthController) Logout(ctx *gin.Context) {
	userID, exist := ctx.Get("user_id")
	if !exist {
		zap.S().Error("error when get user_id from token")
		c.buildErrorResponse(ctx, http.StatusBadRequest, &dtos.GetUserInfoResponse{
			Meta: dtos.BadRequestMeta,
		})
		return
	}
	meta := c.serviceProvider.GetAuthService().Logout(userID.(uint))
	ctx.JSON(meta.Code, meta.Message)

}
func (c *AuthController) Register(ctx *gin.Context) {
	var request *dtos.LoginRequest
	if err := ctx.Bind(&request); err != nil {
		zap.S().Error(ctx, "invalid format register request", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	zap.S().Info(ctx, "Register request with body: ", request)
	meta := c.serviceProvider.GetAuthService().Register(request.Username, request.Password)

	ctx.JSON(meta.Code, meta.Message)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var request *dtos.LoginRequest
	if err := ctx.Bind(&request); err != nil {
		zap.S().Error(ctx, "invalid format request", err)
		c.buildErrorResponse(ctx, http.StatusBadRequest, &dtos.LoginResponse{
			Meta: dtos.BadRequestMeta,
		})
		return
	}

	response := c.serviceProvider.GetAuthService().Login(request.Username, request.Password)
	if response.Meta.Code == http.StatusOK {
		ctx.JSON(http.StatusOK, response.Data)
	} else {
		ctx.JSON(response.Meta.Code, response.Meta.Message)
	}
}

func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var request *dtos.RefreshTokenReq
	if err := ctx.Bind(&request); err != nil {
		zap.S().Error(ctx, "invalid format request", err)
		c.buildErrorResponse(ctx, http.StatusBadRequest, &dtos.LoginResponse{
			Meta: dtos.BadRequestMeta,
		})
		return
	}

	response := c.serviceProvider.GetAuthService().RefreshToken(request.RefreshToken)

	c.buildResponse(ctx, response.Meta.Code, response)
}
