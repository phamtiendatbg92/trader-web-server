package controllers

import (
	"net/http"
	"strconv"
	"trader-web-api/dtos"
	"trader-web-api/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TutorialController struct {
	BaseController
}

func NewTutorialController(sp services.ServiceProvider) *TutorialController {
	c := &TutorialController{}
	c.serviceProvider = sp
	return c
}

/* Upload new post */
func (c *TutorialController) UploadNewPost(ctx *gin.Context) {
	var tutBody dtos.TutorialJson
	ctx.BindJSON(&tutBody)
	result := c.serviceProvider.GetTutService().CreateNewPost(tutBody)
	if result {
		ctx.Status(http.StatusOK)
	} else {
		ctx.Status(http.StatusInternalServerError)
	}
}

func (c *TutorialController) GetAllTutorial(ctx *gin.Context) {
	result := c.serviceProvider.GetTutService().GetAllPost()
	ctx.JSON(200, result)
}
func (c *TutorialController) GetDetailTutorial(ctx *gin.Context) {
	url := ctx.Param("url")
	result := c.serviceProvider.GetTutService().GetDetailTutorial(url)
	if result.Meta.Code == http.StatusOK {
		ctx.JSON(result.Meta.Code, result.Data)
	} else {
		ctx.JSON(result.Meta.Code, result.Meta.Message)
	}
}

func (c *TutorialController) GetAllHashTag(ctx *gin.Context) {
	// tags, _ := dbcontroller.GetHashTag()
	// json.NewEncoder(w).Encode(tags)
	result := c.serviceProvider.GetTutService().GetAllHashTag()
	if result == nil {
		ctx.JSON(500, "")
	} else {
		ctx.JSON(200, result)
	}

}

func (c *TutorialController) UpdateTutorial(ctx *gin.Context) {
	var tutBody dtos.TutorialJson
	ctx.BindJSON(&tutBody)

	result := c.serviceProvider.GetTutService().UpdateTutorial(tutBody)
	ctx.JSON(result.Code, result.Message)
}

func (c *TutorialController) DeleteTutorial(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Wrong ID format")
	} else {
		result := c.serviceProvider.GetTutService().DeleteTutorial(uint(id))
		ctx.JSON(result.Code, result.Message)
	}
}
func (c *TutorialController) PushComment(ctx *gin.Context) {
	userID, exist := ctx.Get("user_id")
	if !exist {
		zap.S().Error("error when get user_id from token")
		c.buildErrorResponse(ctx, http.StatusBadRequest, &dtos.GetUserInfoResponse{
			Meta: dtos.BadRequestMeta,
		})
		return
	}
	var req dtos.Comment1Req
	ctx.Bind(&req)
	result := c.serviceProvider.GetCommentService().PushComment(userID.(uint), req)
	if result != 0 {
		ctx.JSON(http.StatusOK, result)
	} else {
		ctx.JSON(http.StatusInternalServerError, "Cannot save comment")
	}
}
func (c *TutorialController) DeleteComment(ctx *gin.Context) {
	userID, exist := ctx.Get("user_id")
	if !exist {
		zap.S().Error("error when get user_id from token")
		c.buildErrorResponse(ctx, http.StatusBadRequest, &dtos.GetUserInfoResponse{
			Meta: dtos.BadRequestMeta,
		})
		return
	}
	var req dtos.DeleteCommentReq
	error := ctx.Bind(&req)
	if error != nil {
		zap.S().Error("Can not bind body")
	}
	result := c.serviceProvider.GetCommentService().DeleteComment(userID.(uint), req)
	if result {
		ctx.JSON(http.StatusOK, "")
	} else {
		ctx.JSON(http.StatusInternalServerError, "Cannot delete comment")
	}
}
func (c *TutorialController) GetComment(ctx *gin.Context) {
	id := ctx.Query("id")
	temp, _ := strconv.ParseUint(id, 10, 32)
	zap.S().Info("Receive Get All comment request. post_id = ", uint(temp))
	result, meta := c.serviceProvider.GetCommentService().GetAllComment(uint(temp))
	ctx.JSON(meta.Code, result)
}
func (c *TutorialController) PushReply(ctx *gin.Context) {
	userID, exist := ctx.Get("user_id")
	if !exist {
		zap.S().Error("error when get user_id from token")
		c.buildErrorResponse(ctx, http.StatusBadRequest, &dtos.GetUserInfoResponse{
			Meta: dtos.BadRequestMeta,
		})
		return
	}
	var req dtos.Comment2Req
	ctx.Bind(&req)
	id, result := c.serviceProvider.GetCommentService().PushReply(userID.(uint), req)
	if result {
		ctx.JSON(http.StatusOK, id)
	} else {
		ctx.JSON(http.StatusInternalServerError, "Cannot save comment")
	}
}
