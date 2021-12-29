package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ian-antking/king-family-photos/resizePhoto/event"
)

type Handler struct {}

func (h *Handler) Run(_ context.Context, sqsEvent events.SQSEvent) error {
	for _, record := range sqsEvent.Records {
		var message event.Message
		_ = json.Unmarshal([]byte(record.Body), &message)
		fmt.Printf("%+v\n", message)
	}

	return nil
}

func main() {
	handler := Handler{}
	lambda.Start(handler.Run)
}