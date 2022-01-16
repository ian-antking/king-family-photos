package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/ian-antking/king-family-photos/removePhoto/photo"
)

type Handler struct {
	displayBucketName string
	photoRepository   photo.Repository
}

func (h *Handler) Run(_ context.Context, s3Event events.S3Event) error {
	params := h.getPhotoParams(s3Event)

	for _, param := range params {
		err := h.photoRepository.Delete(param)

		if nil != err {
			return err
		}
	}

	return nil
}

func (h *Handler) getPhotoParams(s3Event events.S3Event) []photo.DeletePhotoParams {
	var params []photo.DeletePhotoParams

	for _, record := range s3Event.Records {
		params = append(params, photo.DeletePhotoParams{
			Bucket: h.displayBucketName,
			Key:    record.S3.Object.Key,
		})
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
