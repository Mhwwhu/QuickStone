package router

import (
	"QuickStone/src/web/handlers"
	"QuickStone/src/web/middleware"

	"github.com/gin-gonic/gin"
)

func InitDefaultRouter(router *gin.Engine) {

	userRouter := router.Group("/user")
	{
		//登入注册
		userRouter.POST("/login", handlers.LoginHandle)
		userRouter.POST("/register", handlers.RegisterHandle)
		userRouter.POST("/logout", handlers.LogoutHandle)
	}

	storageRouter := router.Group("/storage", middleware.JwtTokenAuth)
	{
		storageRouter.POST("/upload", handlers.UploadObjectHandle)
		storageRouter.POST("/getbuckets", handlers.GetBucketsHandle)
	}
}
