package photo

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type s3TestSuite struct {
	suite.Suite
	s3 *mockS3Client
}

func (s *s3TestSuite) setUpMocks() {
	s.s3 = new(mockS3Client)
}

func (s *s3TestSuite) TestDelete() {
	s.T().Run("calls DeleteObject with correct input", func(t *testing.T) {
		s.setUpMocks()
		expectedInput := s3.DeleteObjectInput{Bucket: aws.String("bucket"), Key: aws.String("key")}
		s.s3.On("DeleteObject", &expectedInput).Return(&s3.DeleteObjectOutput{}, nil)
		photoRepo := NewS3(s.s3)

		err := photoRepo.Delete(DeletePhotoParams{
			Bucket: "bucket",
			Key:    "key",
		})

		assert.Nil(t, err)
	})
	s.T().Run("relays any errors from s3", func(t *testing.T) {
		s.setUpMocks()
		expectedInput := s3.DeleteObjectInput{Bucket: aws.String("bucket"), Key: aws.String("key")}
		s.s3.On("DeleteObject", &expectedInput).Return(&s3.DeleteObjectOutput{}, errors.New("something went wrong"))
		photoRepo := NewS3(s.s3)

		err := photoRepo.Delete(DeletePhotoParams{
			Bucket: "bucket",
			Key:    "key",
		})

		assert.Equal(t, "error deleting key from bucket: something went wrong", err.Error())
		assert.True(t, errors.Is(err, DeletePhotoError{}))
	})
}

type mockS3Client struct {
	mock.Mock
}

func (m *mockS3Client) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*s3.DeleteObjectOutput), args.Error(1)
}

func TestS3TestSuite(t *testing.T) {
	suite.Run(t, new(s3TestSuite))
}
