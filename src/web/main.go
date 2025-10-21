package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mhwwhu/QuickStone/src/constant/config"
	_ "github.com/mhwwhu/QuickStone/src/web/auth"
	"github.com/sirupsen/logrus"
)

func main() {
	g := gin.Default()
	if err := g.Run(fmt.Sprintf(":%d", config.WebServicePort)); err != nil {
		logrus.Panicf("Can not run GuGoTik Gateway, binding port: %d", config.WebServicePort)
	}
}
