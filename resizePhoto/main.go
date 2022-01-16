package main

import (
	"context"
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

func (h *Handler) getImages(params []photo.GetPhotoParams) ([]photo.GetPhotoOutput, error) {
	var images []photo.GetPhotoOutput
	for _, param := range params {
		getPhotoOutput, err := h.photo.Get(param)

		if nil != err {
			return []photo.GetPhotoOutput{}, err
		}
		images = append(images, getPhotoOutput)
	}

	return images, nil
}

func (h *Handler) processImages(images []photo.GetPhotoOutput) ([]processor.Image, error) {
	var processedImages []processor.Image
	for _, image := range images {
		processedImage, err := h.imageProcessor.Run(processor.Image(image))
		if nil != err {
			return []processor.Image{}, err
		} else {
			processedImages = append(processedImages, processedImage)
		}
	}

	return processedImages, nil
}

func (h *Handler) putImages(images []processor.Image) error {
	for _, image := range images {
		err := h.photo.Put(photo.PutPhotoParams{
			Image:  image.Image,
			Key:    image.Key,
			Bucket: h.displayBucketName,
		})
		if nil != err {
			return err
		}
	}
	return nil
}

func (h *Handler) Run(_ context.Context, s3Event events.S3Event) error {
	params := getPhotoParams(s3Event)

	images, err := h.getImages(params)

	if nil != err {
		return err
	}

	processedImages, err := h.processImages(images)

	if nil != err {
		return err
	}

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
