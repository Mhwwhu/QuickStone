package utils

import (
	"QuickStone/src/common"
	"QuickStone/src/constant"
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

func CreateCtxFromGin(c *gin.Context) context.Context {
	userId, _ := c.Get(constant.CtxUserIdKey)
	userName, _ := c.Get(constant.CtxUserNameKey)
	md := metadata.Pairs(
		constant.CtxUserIdKey, strconv.Itoa(int(userId.(common.UserIdT))),
		constant.CtxUserNameKey, userName.(string),
	)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)
	return ctx
}
