package demo //create demo package for import

import (
	"context" //controls timeouts and cancels
	"fmt"     //i/o stuff
	"time"

	"github.com/dennis-yeom/batman/internal/aws/s3"
	"github.com/dennis-yeom/batman/internal/aws/sqs"
	"github.com/dennis-yeom/batman/internal/redis" //imports redis package
)

// the demo object contains clients for each of the services
type Demo struct {
	redis *redis.RedisClient
	s3    *s3.S3Client
	sqs   *sqs.SQSClient
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

// WithSQS is an option to initialize the SQS client in Demo
func WithSQS(sqsUrl string) DemoOption {
	return func(d *Demo) error {
		sqsClient, err := sqs.NewSQSClient(context.Background(), sqsUrl)
		if err != nil {
			return err
		}
		d.sqs = sqsClient
		return nil
	}
}

// SendMessage sends a message to the SQS queue through the Demo instance
func (d *Demo) SendMessage(ctx context.Context, messageBody string) error {
	if d.sqs == nil {
		return fmt.Errorf("SQS client is not initialized")
	}

	err := d.sqs.SendMessage(ctx, messageBody)
	if err != nil {
		return fmt.Errorf("failed to send message through Demo: %v", err)
	}

	//fmt.Printf("Message successfully sent through Demo: %s\n", messageBody)
	return nil
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

// ListObjectVersions prints all objects in the S3 bucket with their version IDs
func (d *Demo) ListObjectVersions() error {
	// Retrieve all objects with their version IDs
	objects, err := d.s3.GetAllObjectVersions(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list object versions: %w", err)
	}

	// Print each object's key and version ID
	fmt.Println("Objects in bucket with their version IDs:")
	for _, obj := range objects {
		fmt.Printf(" - Key: %s, Version ID: %s\n", obj.Key, obj.VersionID)
	}

	return nil
}

// Watch periodically lists all object versions from S3, compares each version to the Redis cache,
// and updates Redis if necessary. If no changes are detected, it prints a message.
func (d *Demo) Watch(t int) {
	ctx := context.Background()
	// Set up a ticker to run every 30 seconds
	ticker := time.NewTicker(time.Duration(t) * time.Second)
	defer ticker.Stop() // Ensure the ticker is stopped when Watch exits

	fmt.Println("Starting periodic check on S3 object versions...")

	// Using for range to iterate over each tick from ticker.C
	for range ticker.C {
		// Call checkObjectVersions on each tick
		fmt.Println("Ticker triggered: Checking object versions...")
		ov, err := d.s3.GetAllObjectVersions(ctx)
		if err != nil {
			fmt.Printf("Error checking object versions: %v\n", err)
			continue
		}

		//checking for changes
		changesDetected := false

		//for every object from: ov, err := d.s3.GetAllObjectVersions(ctx)
		for _, o := range ov {
			// Track if any changes were detected
			changesDetected := false
			// Retrieve the cached version from Redis
			cachedVersion, err := d.redis.Get(ctx, o.Key)
			if err != nil && err.Error() != "redis: nil" {
				// If there's an error other than key not found, log and continue
				fmt.Printf("Failed to retrieve cache for key %s: %v\n", o.Key, err)
				continue
			}

			// If the key doesn't exist in Redis, set it and print a message
			if cachedVersion == "" {
				if err := d.redis.Set(ctx, o.Key, o.VersionID, 0); err != nil {
					fmt.Printf("Failed to set Redis cache for key %s: %v\n", o.Key, err)
				} else {
					fmt.Printf("Added to cache: Key: %s, Version: %s\n", o.Key, o.VersionID)
					changesDetected = true
					cachedVersion = o.VersionID
				}
			}

			// If the cached version differs from the S3 version, update Redis and print a message
			if cachedVersion != o.VersionID {
				if err := d.redis.Set(ctx, o.Key, o.VersionID, 0); err != nil {
					fmt.Printf("Failed to update Redis cache for key %s: %v\n", o.Key, err)
				} else {
					fmt.Printf("Updated cache: Key: %s, Old Version: %s, New Version: %s\n", o.Key, cachedVersion, o.VersionID)
					changesDetected = true
				}
			}

			if changesDetected {
				fmt.Println("Change detected!")
				d.sqs.SendMessage(ctx, "changes detected!")
			}
		}

		// Print message if no changes were detected during this tick
		if !changesDetected {
			fmt.Println("No changes detected in S3 object versions.")
		}
	}
}
