package routers

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"trader-web-api/conf"
	"trader-web-api/controllers"
	"trader-web-api/dtos"
	"trader-web-api/repositories"
	"trader-web-api/services"
	"trader-web-api/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	TokenTypeJWTAuthen = "Bearer"
)

var jwtHelper utils.JWTHelper

func InitRouter() (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// must use middleware first
	r.Use(middlewareCORS())
	serviceProvider, err := createServiceProvider()
	if err != nil {
		zap.S().Error("Error when createServiceProvider, detail: ", err)
		return nil, err
	}

	authController := controllers.NewAuthController(serviceProvider)
	userController := controllers.NewUserController(serviceProvider)
	tutController := controllers.NewTutorialController(serviceProvider)
	v1 := r.Group("/api/v1/auth")
	v1.POST("/register", authController.Register)
	v1.POST("/login", authController.Login)
	v1.POST("/refresh", authController.RefreshToken)
	v1.Use(middlewareJWTAuthen()).DELETE("/logout", authController.Logout)
	tutRouter := r.Group("/api/v1/tutorial")
	tutRouter.GET("/get-list-tutorials", tutController.GetAllTutorial)
	tutRouter.GET("/detail-tutorial/:url", tutController.GetDetailTutorial)
	tutRouter.GET("/get-hashtag", tutController.GetAllHashTag)
	tutRouter.GET("/comment", tutController.GetComment)
	tutRouter.Use(middlewareJWTAuthen()).POST("/comment", tutController.PushComment)
	tutRouter.Use(middlewareJWTAuthen()).DELETE("/comment", tutController.DeleteComment)
	tutRouter.Use(middlewareJWTAuthen()).POST("/reply-comment", tutController.PushReply)
	tutRouter.Use(middlewareJWTAuthen()).POST("/upload-new-post", tutController.UploadNewPost)
	tutRouter.Use(middlewareJWTAuthen()).PUT("/update-post", tutController.UpdateTutorial)
	tutRouter.Use(middlewareJWTAuthen()).DELETE("/delete-post/:id", tutController.DeleteTutorial)

	v1.Use(middlewareJWTAuthen()).GET("/user", userController.GetInfo)
	return r, nil
}

func middlewareCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		if c.Request.Method == "OPTIONS" {
			c.Writer.WriteHeader(http.StatusOK)
			c.Abort()
			return
		}
		c.Next()
	}
}

func createServiceProvider() (services.ServiceProvider, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True", conf.EnvConfig.MySQL.User, conf.EnvConfig.MySQL.Password, conf.EnvConfig.MySQL.Host, conf.EnvConfig.MySQL.Port, conf.EnvConfig.MySQL.DB)
	dborm, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db, _ := dborm.DB()
	db.SetConnMaxIdleTime(5 * time.Minute)

	// if conf.EnvConfig.Environment == conf.EnvironmentLocal {
	// 	dborm.LogMode(true)
	// }

	repoProvider, err := repositories.NewRepositoryProvider(dborm)
	if err != nil {
		return nil, err
	}

	jHelper, err := utils.NewJWTHelper(conf.EnvConfig.JWT.PublicKey, conf.EnvConfig.JWT.PrivateKey)
	if err != nil {
		return nil, err
	}
	jwtHelper = jHelper

	serviceProvider := services.NewServiceProvider(repoProvider, jHelper)

	return serviceProvider, nil
}

func middlewareJWTAuthen() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		tokens := strings.Split(token, " ")
		if len(tokens) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &dtos.BaseResponse{
				Meta: &dtos.Meta{
					Code:    http.StatusUnauthorized,
					Message: "Unauthorized"}})
			return
		}

		switch tokens[0] {
		case TokenTypeJWTAuthen:
			var userClaim dtos.AuthClaims
			err := jwtHelper.ParseClaims(tokens[1], &userClaim)
			if err != nil {
				zap.S().Error("Error when parse jwt token, detail: ", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, &dtos.BaseResponse{
					Meta: &dtos.Meta{
						Code:    http.StatusUnauthorized,
						Message: "Unauthorized"}})
				return
			}
			c.Set("user_id", userClaim.UserID)
		default:
			c.AbortWithStatusJSON(http.StatusUnauthorized, &dtos.BaseResponse{
				Meta: &dtos.Meta{
					Code:    http.StatusUnauthorized,
					Message: "Unauthorized"}})
			return
		}
		c.Next()
	}
}
