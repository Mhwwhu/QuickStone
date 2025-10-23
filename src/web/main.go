package main

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/mhwwhu/QuickStone/src/config"
	routers "github.com/mhwwhu/QuickStone/src/web/router"
	"github.com/sirupsen/logrus"
)

func main() {
	g := gin.Default()

	//先注册中间件
	store := cookie.NewStore([]byte("secret"))
	g.Use(sessions.Sessions("mysession", store))

	routers.InitDefaultRouter(g)

	if err := g.Run(fmt.Sprintf(":%d", config.WebServicePort)); err != nil {
		logrus.Panicf("Cannot run gateway, binding port: %d", config.WebServicePort)
	}
}
