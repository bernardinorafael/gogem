package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Client struct {
	sqs      *sqs.Client
	queueURL string
}

type Message struct {
	ID            string
	Body          string
	ReceiptHandle string
}

type Config struct {
	Region   string
	QueueURL string
}

func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("sqs: failed to load AWS config: %w", err)
	}

	return &Client{
		sqs:      sqs.NewFromConfig(awsCfg),
		queueURL: cfg.QueueURL,
	}, nil
}

func (c *Client) Publish(ctx context.Context, message any) error {
	b, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("sqs: failed to marshal JSON: %w", err)
	}

	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(c.queueURL),
		MessageBody: aws.String(string(b)),
	}

	if _, err := c.sqs.SendMessage(ctx, input); err != nil {
		return fmt.Errorf("sqs: failed to publish message: %w", err)
	}

	return nil
}

func (c *Client) Consume(ctx context.Context) ([]Message, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueURL),
		MaxNumberOfMessages: *aws.Int32(10),
		WaitTimeSeconds:     *aws.Int32(20),
	}

	result, err := c.sqs.ReceiveMessage(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("sqs: failed to consume messages: %w", err)
	}

	messages := make([]Message, 0, len(result.Messages))
	for _, msg := range result.Messages {
		m := Message{
			ID:            aws.ToString(msg.MessageId),
			Body:          aws.ToString(msg.Body),
			ReceiptHandle: aws.ToString(msg.ReceiptHandle),
		}
		messages = append(messages, m)
	}

	return messages, nil
}

func (c *Client) DeleteMessage(ctx context.Context, receiptHandle string) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	}

	_, err := c.sqs.DeleteMessage(ctx, input)
	if err != nil {
		return fmt.Errorf("sqs: failed to delete message: %w", err)
	}

	return nil
}
