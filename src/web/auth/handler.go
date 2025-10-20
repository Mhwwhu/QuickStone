package auth

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mhwwhu/QuickStone/src/constant/config"
	"github.com/mhwwhu/QuickStone/src/rpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Client auth.AuthServiceClient

func init() {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", config.EnvCfg.PodIpAddr, config.Config.AuthServerPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Fail to connect: %v\n", err)
	}

	defer conn.Close()
}

func LoginHandle(c *gin.Context) {

}
