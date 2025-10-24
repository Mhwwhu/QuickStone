package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mhwwhu/QuickStone/src/config"
	"github.com/mhwwhu/QuickStone/src/web/midware"
	routers "github.com/mhwwhu/QuickStone/src/web/router"
	"github.com/sirupsen/logrus"
)

func main() {
	g := gin.Default()

	//先注册中间件
	midware.InitSession(g)

	routers.InitDefaultRouter(g)

	if err := g.Run(fmt.Sprintf(":%d", config.WebServicePort)); err != nil {
		logrus.Panicf("Cannot run gateway, binding port: %d", config.WebServicePort)
	}

}
