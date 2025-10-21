package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mhwwhu/QuickStone/src/config"
	"github.com/mhwwhu/QuickStone/src/web/auth"
	"github.com/sirupsen/logrus"
)

func main() {
	g := gin.Default()

	rootPath := g.Group("/")

	user := rootPath.Group("/user")
	user.POST("/login", auth.LoginHandle)
	user.POST("/register", auth.RegisterHandle)

	if err := g.Run(fmt.Sprintf(":%d", config.WebServicePort)); err != nil {
		logrus.Panicf("Cannot run gateway, binding port: %d", config.WebServicePort)
	}
}
