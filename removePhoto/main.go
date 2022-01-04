package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ian-antking/king-family-photos/removePhoto/event"
	"github.com/ian-antking/king-family-photos/removePhoto/photo"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Handler struct {
	displayBucketName string
	photoRepository   photo.Repository
}

func (h *Handler) Run(_ context.Context, sqsEvent events.SQSEvent) error {
	messages := getMessages(sqsEvent)
	params := h.getPhotoParams(messages)

	for _, param := range params {
		err := h.photoRepository.Delete(param)

		if nil != err {
			log.Fatalf("Error deleting image from s3: %s", err.Error())
		}
	}

	return nil
}

func getMessages(sqsEvent events.SQSEvent) []event.Message {
	var messages []event.Message
	for _, record := range sqsEvent.Records {
		var message event.Message
		_ = json.Unmarshal([]byte(record.Body), &message)
		messages = append(messages, message)
	}

	return messages
}

func (h *Handler) getPhotoParams(messages []event.Message) []photo.DeletePhotoParams {
	var params []photo.DeletePhotoParams

	for _, message := range messages {
		for _, record := range message.Records {
			params = append(params, photo.DeletePhotoParams{
				Bucket: h.displayBucketName,
				Key:    record.S3.Object.Key,
			})
		}
	}

	return params
}

func NewHandler(bucketName string, repository photo.Repository) Handler {
	return Handler{
		displayBucketName: bucketName,
		photoRepository:   repository,
	}
}

func main() {
	displayBucketName := os.Getenv("DISPLAY_BUCKET")

	awsSession := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		},
	))

	s3Client := s3.New(awsSession)
	photoRepository := photo.NewS3(s3Client)

	handler := NewHandler(displayBucketName, &photoRepository)

	lambda.Start(handler.Run)
}
