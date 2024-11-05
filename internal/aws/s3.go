package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 interface {
	//CheckConnection(ctx context.Context) error // Method to check S3 connectivity
	ListFiles(ctx context.Context) error
	GetObjectVersion(ctx context.Context, key string) (string, error)
	GetAllObjectVersions(ctx context.Context) ([]ObjectInfo, error)
}

type S3Client struct {
	client *s3.Client
	bucket string
}

// NewS3Client initializes a new S3Client with the specified bucket
func NewS3Client(ctx context.Context, bucket string) (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1")) //make sure region is correct
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return &S3Client{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
	}, nil
}

// CheckConnection tests connectivity to the S3 bucket by attempting to list objects
// func (s *S3Client) CheckConnection(ctx context.Context) error {
// 	input := &s3.ListObjectsV2Input{
// 		Bucket: aws.String(s.bucket),
// 	}

// 	_, err := s.client.ListObjectsV2(ctx, input)
// 	if err != nil {
// 		return fmt.Errorf("failed to connect to S3 bucket %s: %w", s.bucket, err)
// 	}

// 	fmt.Printf("Successfully connected to S3 bucket: %s\n", s.bucket)
// 	return nil
// }

// GetObjectVersion retrieves the metadata of an object and returns its version ID
func (s *S3Client) GetObjectVersion(ctx context.Context, key string) (string, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	result, err := s.client.HeadObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to get object metadata: %w", err)
	}

	versionID := aws.ToString(result.VersionId)
	//fmt.Printf("File %s in bucket %s has version ID: %s\n", key, s.bucket, versionID)
	return versionID, nil
}

// ListFiles lists all files in the S3 bucket and prints their version IDs
func (s *S3Client) ListFiles(ctx context.Context) error {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	}

	// List all objects in the specified bucket
	result, err := s.client.ListObjectsV2(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	fmt.Println("Files in bucket:", s.bucket)
	for _, item := range result.Contents {
		// Use GetObjectVersion to retrieve the version ID for each object
		versionID, err := s.GetObjectVersion(ctx, *item.Key)
		if err != nil {
			// Print an error if we can't retrieve the version, but continue with the next item
			fmt.Printf(" - %s (size: %d) - failed to retrieve version: %v\n", *item.Key, item.Size, err)
			continue
		}

		// Print the key, size, and version ID of the object
		fmt.Printf(" - %s (size: %d, version: %s)\n", *item.Key, item.Size, versionID)
	}
	return nil
}

// ObjectInfo holds the key (filename) and version ID of an object
type ObjectInfo struct {
	Key       string
	VersionID string
}

// GetAllObjectVersions retrieves the filename and version ID for all objects in the S3 bucket
func (s *S3Client) GetAllObjectVersions(ctx context.Context) ([]ObjectInfo, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	}

	result, err := s.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	var objects []ObjectInfo
	for _, item := range result.Contents {
		versionID, err := s.GetObjectVersion(ctx, *item.Key)
		if err != nil {
			fmt.Printf("Failed to get version for object %s: %v\n", *item.Key, err)
			continue
		}

		objects = append(objects, ObjectInfo{
			Key:       *item.Key,
			VersionID: versionID,
		})
	}

	return objects, nil
}

var _ S3 = &S3Client{}
