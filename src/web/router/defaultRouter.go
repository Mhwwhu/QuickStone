package router

import (
	"QuickStone/src/web/handlers/bucket"
	trans "QuickStone/src/web/handlers/transmission"
	"QuickStone/src/web/handlers/user"
	"QuickStone/src/web/middleware"

	"github.com/gin-gonic/gin"
)

func InitDefaultRouter(router *gin.Engine) {

	userRouter := router.Group("/user")
	{
		//登入注册
		userRouter.POST("/login", user.LoginHandle)
		userRouter.POST("/register", user.RegisterHandle)
		userRouter.GET("/checklogin", user.CheckLoginHandle)
	}

	storageRouter := router.Group("/storage", middleware.JwtTokenAuth)
	{
		storageRouter.POST("/upload", trans.UploadObjectHandle)
	}

	bucketRouter := router.Group("/bucket", middleware.JwtTokenAuth)
	{
		bucketRouter.POST("/create", bucket.CreateBucketHandle)
		bucketRouter.GET("/info", bucket.ShowBucketHandle)
	}
}
