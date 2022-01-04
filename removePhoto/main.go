package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Handler struct {
	displayBucketName string
}


func (h *Handler) Run(_ context.Context, sqsEvent events.SQSEvent) error {
	fmt.Printf("%+v\n", sqsEvent)
	return nil
}

func NewHandler(bucketName string) Handler {
	return Handler{
		displayBucketName: bucketName,
	}
}

func main() {
	displayBucketName := os.Getenv("DISPLAY_BUCKET")

	handler := NewHandler(displayBucketName)

	lambda.Start(handler.Run)
}