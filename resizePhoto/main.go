package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/ian-antking/king-family-photos/resizePhoto/photo"
	"github.com/ian-antking/king-family-photos/resizePhoto/processor"
)

type Handler struct {
	photo             photo.Repository
	displayBucketName string
	imageProcessor    processor.Processor
}

func (h *Handler) getImages(params []photo.GetPhotoParams) ([]photo.GetPhotoOutput, error) {
	var images []photo.GetPhotoOutput
	for _, param := range params {
		getPhotoOutput, err := h.photo.Get(param)

		if nil != err {
			return []photo.GetPhotoOutput{}, fmt.Errorf("error getting %s from %s: %s", param.Key, param.Bucket, err.Error())
		}
		images = append(images, getPhotoOutput)
	}

	return images, nil
}

func (h *Handler) putImage(image processor.Image) error {
	err := h.photo.Put(photo.PutPhotoParams{
		Image:  image.Image,
		Key:    image.Key,
		Bucket: h.displayBucketName,
	})

	return err
}

func (h *Handler) processImages(images []photo.GetPhotoOutput) []processor.Image {
	var processedImages []processor.Image
	for _, image := range images {
		processedImage, err := h.imageProcessor.Run(processor.Image(image))
		if nil != err {
			log.Fatalf("Error processing image: %s", err.Error())
		}
		processedImages = append(processedImages, processedImage)
	}

	return processedImages
}

func (h *Handler) putImages(images []processor.Image) error {
	for _, image := range images {
		err := h.putImage(image)
		if nil != err {
			return fmt.Errorf("error getting %s from %s: %s", image.Key, image.Bucket, err.Error())
		}
	}
	return nil
}

func getPhotoParams(s3Event events.S3Event) []photo.GetPhotoParams {
	var params []photo.GetPhotoParams

	for _, message := range s3Event.Records {
		params = append(params, photo.GetPhotoParams{
			Bucket: message.S3.Bucket.Name,
			Key:    message.S3.Object.Key,
		})
	}

	return params
}

func (h *Handler) Run(_ context.Context, s3Event events.S3Event) error {
	params := getPhotoParams(s3Event)

	images, err := h.getImages(params)

	if nil != err {
		return err
	}

	processedImages := h.processImages(images)
	err = h.putImages(processedImages)

	return err
}

func NewHandler(repository photo.Repository, bucketName string, imageProcessor processor.Processor) Handler {
	return Handler{
		photo:             repository,
		displayBucketName: bucketName,
		imageProcessor:    imageProcessor,
	}
}

func main() {
	displayBucketName := os.Getenv("DISPLAY_BUCKET")

	awsSession := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		},
	))

	s3Downloader := s3manager.NewDownloader(awsSession)
	s3Uploader := s3manager.NewUploader(awsSession)
	photoRepository := photo.NewS3(s3Downloader, s3Uploader)
	imageProcessor := processor.NewResizer(0, 480)
	handler := NewHandler(&photoRepository, displayBucketName, &imageProcessor)

	lambda.Start(handler.Run)
}
