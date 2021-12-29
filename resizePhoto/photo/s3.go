package photo

import (
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Client interface {
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

type S3 struct {
	client s3Client
}

func (s *S3) Get(params GetPhotoParams) (GetPhotoOutput, error) {
	return GetPhotoOutput{}, nil
}

func NewS3(client s3Client) S3 {
	return S3{
		client: client,
	}
}
