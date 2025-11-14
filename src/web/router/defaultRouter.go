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
		userRouter.POST("/logout", user.LogoutHandle)
	}

	storageRouter := router.Group("/storage", middleware.JwtTokenAuth)
	{
		storageRouter.POST("/upload", trans.UploadObjectHandle)

		bucketRouter := storageRouter.Group("/bucket")
		{
			bucketRouter.POST("/create", bucket.CreateBucketHandle)
			bucketRouter.GET("/info", bucket.ShowBucketHandle)
			bucketRouter.POST("/overview", bucket.ShowUserBucketsHandle)
			bucketRouter.GET("/objects", bucket.ShowObjectsHandle)
		}
	}
}
