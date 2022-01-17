package photo

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type s3TestSuite struct {
	suite.Suite
	downloader *mockS3Downloader
	uploader   *mockS3Uploader
}

func (s *s3TestSuite) setUpMocks() {
	s.downloader = new(mockS3Downloader)
	s.uploader = new(mockS3Uploader)
}

func (s *s3TestSuite) TestGet() {
	s.T().Run("calls Download with correct input", func(t *testing.T) {
		s.setUpMocks()
		photoRepo := NewS3(s.downloader, s.uploader)

		s.downloader.On(
			"Download",
			&aws.WriteAtBuffer{},
			&s3.GetObjectInput{
				Bucket: aws.String("ingestBucket"),
				Key:    aws.String("photoKey"),
			},
			mock.Anything,
		).Return(nil)

		data := []byte("data")
		idx := int64(len(data))
		buffer := &aws.WriteAtBuffer{}
		_, _ = buffer.WriteAt(data, idx)

		expected := GetPhotoOutput{
			Image:  buffer.Bytes(),
			Bucket: "ingestBucket",
			Key:    "photoKey",
		}

		actual, err := photoRepo.Get(GetPhotoParams{
			Bucket: "ingestBucket",
			Key:    "photoKey",
		})

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	s.T().Run("forwards errors return from s3", func(t *testing.T) {
		s.setUpMocks()
		photoRepo := NewS3(s.downloader, s.uploader)

		s.downloader.On(
			"Download",
			&aws.WriteAtBuffer{},
			&s3.GetObjectInput{
				Bucket: aws.String("ingestBucket"),
				Key:    aws.String("photoKey"),
			},
			mock.Anything,
		).Return(errors.New("something went wrong"))

		_, err := photoRepo.Get(GetPhotoParams{
			Bucket: "ingestBucket",
			Key:    "photoKey",
		})

		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, GetPhotoError{}))
		assert.Equal(t, "error getting photoKey from ingestBucket: something went wrong", err.Error())
	})
}

func (s *s3TestSuite) TestPut() {
	s.T().Run("calls Upload with correct input", func(t *testing.T) {
		s.setUpMocks()
		photoRepo := NewS3(s.downloader, s.uploader)

		params := PutPhotoParams{
			Image:  []byte{},
			Key:    "photoKey",
			Bucket: "displayBucket",
		}

		s.uploader.On("Upload", &s3manager.UploadInput{
			Body:   bytes.NewReader([]byte{}),
			Bucket: aws.String("displayBucket"),
			Key:    aws.String("photoKey"),
		}, mock.Anything).Return(&s3manager.UploadOutput{}, nil)

		err := photoRepo.Put(params)

		assert.Nil(t, err)
	})

	s.T().Run("forwards errors from s3", func(t *testing.T) {
		s.setUpMocks()
		photoRepo := NewS3(s.downloader, s.uploader)

		params := PutPhotoParams{
			Image:  []byte{},
			Key:    "photoKey",
			Bucket: "displayBucket",
		}

		s.uploader.On("Upload", &s3manager.UploadInput{
			Body:   bytes.NewReader([]byte{}),
			Bucket: aws.String("displayBucket"),
			Key:    aws.String("photoKey"),
		}, mock.Anything).Return(&s3manager.UploadOutput{}, errors.New("something went wrong"))

		err := photoRepo.Put(params)

		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, PutPhotoError{}))
		assert.Equal(t, "error putting photoKey in displayBucket: something went wrong", err.Error())
	})
}

type mockS3Downloader struct {
	mock.Mock
}

func (m *mockS3Downloader) Download(writer io.WriterAt, input *s3.GetObjectInput, f ...func(*s3manager.Downloader)) (int64, error) {
	args := m.Called(writer, input, f)

	data := []byte("data")

	idx := int64(len(data))

	_, _ = writer.WriteAt(data, idx)

	return idx, args.Error(0)
}

type mockS3Uploader struct {
	mock.Mock
}

func (m *mockS3Uploader) Upload(input *s3manager.UploadInput, f ...func(uploader *s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	args := m.Called(input, f)
	return args.Get(0).(*s3manager.UploadOutput), args.Error(1)
}

func TestS3TestSuite(t *testing.T) {
	suite.Run(t, new(s3TestSuite))
}
