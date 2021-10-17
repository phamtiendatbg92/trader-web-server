package controllers

import (
	"trader-web-api/services"

	"github.com/gin-gonic/gin"
)

type BaseController struct {
	serviceProvider services.ServiceProvider
}

func (b *BaseController) buildSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, data)
}

func (b *BaseController) buildErrorResponse(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

func (b *BaseController) buildResponse(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}
