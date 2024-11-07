package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// struct to hold sqs client info
type SQSClient struct {
	client   *sqs.Client
	queueURL string
}

// func creates new sqs client. returns pointer to sqs client struct
func NewSQSClient(ctx context.Context, queueURL string) (*SQSClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	return &SQSClient{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
	}, nil
}

// SendMessage sends a message to the SQS queue
func (s *SQSClient) SendMessage(ctx context.Context, messageBody string) error {
	input := &sqs.SendMessageInput{
		QueueUrl:    &s.queueURL,
		MessageBody: aws.String(messageBody),
	}

	// Send the message
	resp, err := s.client.SendMessage(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	// Log the message ID to confirm success
	fmt.Printf("Message sent to queue %s, Message ID: %s\n", s.queueURL, *resp.MessageId)
	return nil
}
