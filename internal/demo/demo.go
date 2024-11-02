package demo //create demo package for import

import (
	"context" //controls timeouts and cancels
	"fmt"     //i/o stuff

	"github.com/dennis-yeom/batman/internal/redis" //imports redis package
)

// the demo object contains a client for redis
type Demo struct {
	redis *redis.RedisClient
}

type DemoOption func(*Demo) error

// New initializes a new Demo instance with Redis
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

// Set sets a value in Redis
func (d *Demo) Set(key string, value string) error {
	fmt.Println("Running d.Set()")

	ctx := context.Background()

	// Sets the value to the key in Redis
	if err := d.redis.Set(ctx, key, value, 0); err != nil {
		return err
	}

	fmt.Printf("Set value:%s to key:%s\n", value, key)
	return nil
}

// Get retrieves a value from Redis
func (d *Demo) Get(key string) error {
	fmt.Println("Running d.Get()")

	ctx := context.Background()

	// Gets the value of the key in Redis
	val, err := d.redis.Get(ctx, key)
	if err != nil {
		return err
	}

	fmt.Printf("Got value for key %s: %s\n", key, val)
	return nil
}
