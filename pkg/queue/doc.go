// Package queue provides an AWS SQS client wrapper for publishing, consuming,
// and deleting messages from a queue.
//
// Creating a client:
//
//	client, err := queue.NewClient(ctx, queue.Config{
//	    Region:   "us-east-1",
//	    QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789/my-queue",
//	})
//
// Publishing a message (automatically serialized as JSON):
//
//	err := client.Publish(ctx, OrderCreatedEvent{
//	    OrderID: "order_abc123",
//	    Total:   9900,
//	})
//
// Consuming messages (long-polling, up to 10 messages per call):
//
//	messages, err := client.Consume(ctx)
//	for _, msg := range messages {
//	    // msg.ID, msg.Body, msg.ReceiptHandle
//	    processMessage(msg)
//	    client.DeleteMessage(ctx, msg.ReceiptHandle)
//	}
//
// The client uses AWS default credential chain resolution (environment variables,
// shared credentials file, IAM roles, etc.).
package queue
