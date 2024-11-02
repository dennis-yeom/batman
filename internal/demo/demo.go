package demo //create demo package for import

import (
	"context" //controls timeouts and cancels
	"fmt"     //i/o stuff

	s3 "github.com/dennis-yeom/batman/internal/aws"
	"github.com/dennis-yeom/batman/internal/redis" //imports redis package
)

// the demo object contains a client for redis
type Demo struct {
	redis *redis.RedisClient
	s3    *s3.S3Client
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

// WithS3 sets up the S3 client for the Demo struct
func WithS3(bucket string) DemoOption {
	return func(d *Demo) error {
		ctx := context.TODO()
		s3Client, err := s3.NewS3Client(ctx, bucket)
		if err != nil {
			return fmt.Errorf("failed to initialize S3 client: %w", err)
		}
		d.s3 = s3Client // Directly assign the S3Client instance
		return nil
	}
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

func (d *Demo) List() error {
	return d.s3.ListFiles(context.Background())
}
