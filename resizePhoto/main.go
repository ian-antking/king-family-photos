package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ian-antking/king-family-photos/resizePhoto/event"
	"github.com/ian-antking/king-family-photos/resizePhoto/photo"
)

type Handler struct {
	photo             photo.Repository
	displayBucketName string
}

func (h *Handler) getImages(params []photo.GetPhotoParams) []photo.GetPhotoOutput {
	var images []photo.GetPhotoOutput
	for _, param := range params {
			getPhotoOutput, err := h.photo.Get(param)

			if nil != err {
				fmt.Println(err.Error())
			}
		images = append(images, getPhotoOutput)
	}

	return images
}

func (h *Handler) putImage(image photo.GetPhotoOutput) error {
	err := h.photo.Put(photo.PutPhotoParams{
		Image:  image.Image,
		Key:    image.Key,
		Bucket: h.displayBucketName,
	})

	return err
}

func (h *Handler) processImages(images []photo.GetPhotoOutput) {
	for _, image := range images {
		err := h.putImage(image)
		if nil != err {
			fmt.Println(err.Error())
		}
	}
}

func (h *Handler) Run(_ context.Context, sqsEvent events.SQSEvent) error {
	messages := getMessages(sqsEvent)
	params := getPhotoParams(messages)


	images := h.getImages(params)

	h.processImages(images)

	return nil
}

func NewHandler(repository photo.Repository, bucketName string) Handler {
	return Handler{
		photo:             repository,
		displayBucketName: bucketName,
	}
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

func getPhotoParams(messages []event.Message) []photo.GetPhotoParams {
	var params []photo.GetPhotoParams

	for _, message := range messages {
		for _, record := range message.Records {
			params = append(params, photo.GetPhotoParams{
				Bucket: record.S3.Bucket.Name,
				Key:    record.S3.Object.Key,
			})
		}
	}

	return params
}

func main() {
	displayBucketName := os.Getenv("DISPLAY_BUCKET")

	awsSession := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		},
	))

	s3Client := s3.New(awsSession)
	s3Downloader := s3manager.NewDownloader(awsSession)
	s3Uploader := s3manager.NewUploader(awsSession)
	photoRepository := photo.NewS3(s3Client, s3Downloader, s3Uploader)
	handler := NewHandler(&photoRepository, displayBucketName)

	lambda.Start(handler.Run)
}
