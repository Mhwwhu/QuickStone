package midware

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"

	"github.com/gin-gonic/gin"
)

func InitSession(g *gin.Engine) {

	store, err := redis.NewStore(10, "tcp", "localhost:6379", "", "", []byte("secret"))
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	g.Use(sessions.Sessions("mysession", store))
}
