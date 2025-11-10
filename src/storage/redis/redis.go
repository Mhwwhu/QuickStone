package redis

import (
	"QuickStone/src/config"
	"strings"

	"github.com/redis/go-redis/v9"
)

var Client redis.UniversalClient

func init() {
	addrs := strings.Split(config.EnvCfg.RedisAddr, ";")
	Client = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:      addrs,
		Password:   config.EnvCfg.RedisPassword,
		DB:         config.EnvCfg.RedisDB,
		MasterName: config.EnvCfg.RedisMaster,
	})
}
