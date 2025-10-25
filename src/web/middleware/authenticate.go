package middleware

import "github.com/gin-gonic/gin"

func Authenticate(c *gin.Context) {
	userID := 1
	c.Set("UserId", userID)
	c.Next()
}
