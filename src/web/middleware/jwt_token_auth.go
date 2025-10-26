package middleware

import (
	"QuickStone/src/constant"
	"QuickStone/src/utils/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func JwtTokenAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := jwt.VerifyToken(tokenString)
	if err != nil {
		logrus.Infof("Token is not valid: token = %s, err = %v", tokenString, err)
		c.JSON(http.StatusOK, gin.H{
			constant.StatusCodeKey: constant.UnauthorizedErrorCode,
			constant.StatusMsgKey:  constant.UnauthorizedError,
		})
		c.Abort()
	}
	c.Set(constant.CtxClaimKey, claims)
	c.Next()
}
