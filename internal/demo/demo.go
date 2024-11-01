package demo

import (
	"github.com/dennis-yeom/batman/internal/redis"
)

type Demo struct {
	redis *redis.RedisClient
}

type DemoOption func(*Demo) error

// New initializes a new Demo instance with Redis and optional S3 configuration
func New(port int, opts ...DemoOption) (*Demo, error) {
	d := &Demo{
		redis: redis.New(port),
	}

	for _, opt := range opts {
		if err := opt(d); err != nil {
			return nil, err
		}
	}

	return d, nil
}
