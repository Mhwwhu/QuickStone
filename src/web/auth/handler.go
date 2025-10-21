package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/mhwwhu/QuickStone/src/constant/config"
	"github.com/mhwwhu/QuickStone/src/rpc/auth"
	grpcutil "github.com/mhwwhu/QuickStone/src/utils/grpc"
)

var Client auth.AuthServiceClient

func init() {
	conn := grpcutil.Connect(config.AuthServerName)
	Client = auth.NewAuthServiceClient(conn)
}

func LoginHandle(c *gin.Context) {

}
