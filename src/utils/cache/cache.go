package cache

import (
	"QuickStone/src/config"
	"QuickStone/src/storage/redis"
	"context"
	"encoding/json"
)

// var c *cache.Cache

func init() {
	// c = cache.New(5*time.Minute, 10*time.Minute)
}

func Set(ctx context.Context, key string, obj interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	result := redis.Client.Set(ctx, key, b, config.RedisExpirationTime)
	return result.Err()
}

func Get[T any](ctx context.Context, key string) *T {
	obj, err := redis.Client.Get(ctx, key).Result()
	if err != nil {
		return nil
	}
	var v T
	if err := json.Unmarshal([]byte(obj), &v); err != nil {
		return nil
	}
	return &v
}
