package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Handler struct {}

func (h *Handler) Run(_ context.Context, event events.SQSEvent) error {
	for _, record := range event.Records {
		fmt.Println(record.Body)
	}

	return nil
}

func main() {
	handler := Handler{}
	lambda.Start(handler.Run)
}