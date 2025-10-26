package main

import (
	"fmt"

	"QuickStone/src/config"
	"QuickStone/src/web/router"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	g := gin.Default()

	router.InitDefaultRouter(g)

	if err := g.Run(fmt.Sprintf(":%d", config.WebServicePort)); err != nil {
		logrus.Panicf("Cannot run gateway, binding port: %d", config.WebServicePort)
	}

}
