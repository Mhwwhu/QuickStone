package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/mhwwhu/QuickStone/src/web/controller"
)

func InitDefaultRouter(router *gin.Engine) {

	userRouter := router.Group("/user")
	{
		//登入注册
		userRouter.POST("/login", controller.LoginHandle)
		userRouter.POST("/register", controller.RegisterHandle)
		userRouter.GET("/checklogin", controller.CheckLogin)
		//数据传输
		userRouter.POST("/upload", controller.UploadFile)
	}

}
