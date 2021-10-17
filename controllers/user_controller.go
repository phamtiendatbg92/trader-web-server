package controllers

import (
	"net/http"
	"trader-web-api/dtos"
	"trader-web-api/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	BaseController
}

func NewUserController(sp services.ServiceProvider) *UserController {
	c := &UserController{}
	c.serviceProvider = sp
	return c
}

func (c *UserController) GetInfo(ctx *gin.Context) {
	userID, exist := ctx.Get("user_id")
	if !exist {
		zap.S().Error("error when get user_id from token")
		c.buildErrorResponse(ctx, http.StatusBadRequest, &dtos.GetUserInfoResponse{
			Meta: dtos.BadRequestMeta,
		})
		return
	}
	res := c.serviceProvider.GetUserService().GetUserInfo(userID.(uint))
	if res.Meta.Code == http.StatusOK {
		ctx.JSON(http.StatusOK, res)
	} else {
		ctx.JSON(res.Meta.Code, res.Meta.Message)
	}
}
