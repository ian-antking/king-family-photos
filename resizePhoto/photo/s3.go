package photo

import (
	"bytes"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type s3Downloader interface {
	Download(io.WriterAt, *s3.GetObjectInput, ...func(*s3manager.Downloader)) (int64, error)
}

type s3Uploader interface {
	Upload(*s3manager.UploadInput, ...func(uploader *s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

type S3 struct {
	downloader s3Downloader
	uploader   s3Uploader
}

func (s *S3) Get(params GetPhotoParams) (GetPhotoOutput, error) {
	getObjectInput := s3.GetObjectInput{
		Bucket: aws.String(params.Bucket),
		Key:    aws.String(params.Key),
	}

	buffer := &aws.WriteAtBuffer{}

	_, err := s.downloader.Download(buffer, &getObjectInput)

	if nil != err {
		return GetPhotoOutput{}, GetPhotoError{Err: fmt.Errorf("error getting %s from %s: %s", params.Key, params.Bucket, err.Error())}
	}

	output := GetPhotoOutput{
		Bucket: params.Bucket,
		Key:    params.Key,
		Image:  buffer.Bytes(),
	}

	return output, nil
}

func (s *S3) Put(params PutPhotoParams) error {
	putObjectInput := s3manager.UploadInput{
		Body:   bytes.NewReader(params.Image),
		Bucket: aws.String(params.Bucket),
		Key:    aws.String(params.Key),
	}

	_, err := s.uploader.Upload(&putObjectInput)

	if nil != err {
		return PutPhotoError{Err: fmt.Errorf("error putting %s in %s: %s", params.Key, params.Bucket, err.Error())}
	}

	return nil
}

func NewS3(downloader s3Downloader, uploader s3Uploader) S3 {
	return S3{
		downloader: downloader,
		uploader:   uploader,
	}
}
