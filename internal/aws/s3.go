package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 interface {
	CheckConnection(ctx context.Context) error // Method to check S3 connectivity
}

type S3Client struct {
	client *s3.Client
	bucket string
}

// NewS3Client initializes a new S3Client with the specified bucket
func NewS3Client(ctx context.Context, bucket string) (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-2"))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return &S3Client{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
	}, nil
}

// CheckConnection tests connectivity to the S3 bucket by attempting to list objects
func (s *S3Client) CheckConnection(ctx context.Context) error {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	}

	_, err := s.client.ListObjectsV2(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to connect to S3 bucket %s: %w", s.bucket, err)
	}

	fmt.Printf("Successfully connected to S3 bucket: %s\n", s.bucket)
	return nil
}

// ListFiles lists all files in the S3 bucket
func (s *S3Client) ListFiles(ctx context.Context) error {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	}

	result, err := s.client.ListObjectsV2(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	fmt.Println("Files in bucket:", s.bucket)
	for _, item := range result.Contents {
		fmt.Printf(" - %s (size: %d)\n", *item.Key, item.Size)
	}
	return nil
}

var _ S3 = &S3Client{}
