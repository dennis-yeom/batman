package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RedisClient defines a redis client.
type RedisClient struct {
	client *redis.Client
}

// New returns a new RedisClient
func New(port int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("localhost:%d", port),
	})

	return &RedisClient{
		client: rdb,
	}
}
