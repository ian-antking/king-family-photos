package photo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Client interface {
	DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
}

type S3 struct {
	client        s3Client
	displayBucket string
}

func (s *S3) Delete(params DeletePhotoParams) error {
	deletePhotoInput := s3.DeleteObjectInput{
		Bucket: aws.String(params.Bucket),
		Key:    aws.String(params.Key),
	}

	_, err := s.client.DeleteObject(&deletePhotoInput)

	return err
}

func NewS3(client s3Client) S3 {
	return S3{
		client: client,
	}
}
