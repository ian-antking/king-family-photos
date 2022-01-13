package main

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/ian-antking/king-family-photos/removePhoto/photo"
)

type handlerTestSuite struct {
	suite.Suite
	photoRepository *mockPhotoRepository
}

func (s *handlerTestSuite) TestGetPhotoParams() {
	s.T().Run("converts records on an S3Event to DeletePhotoParams", func(t *testing.T) {
		s.setupMocks()
		handler := NewHandler("displayBucket", s.photoRepository)
		event := events.S3Event{
			Records: []events.S3EventRecord{
				{
					S3: events.S3Entity{
						Bucket: events.S3Bucket{
							Name: "ingestBucket",
						},
						Object: events.S3Object{
							Key: "photoKey",
						},
					},
				},
			},
		}
		expected := []photo.DeletePhotoParams{
			{
				Bucket: "displayBucket",
				Key:    "photoKey",
			},
		}

		result := handler.getPhotoParams(event)

		assert.Equal(t, expected, result)
	})
}

func (s *handlerTestSuite) TestRun() {
	s.T().Run("processes s3 event and deletes photos from display bucket", func(t *testing.T) {
		s.setupMocks()
		handler := NewHandler("displayBucket", s.photoRepository)
		event := events.S3Event{
			Records: []events.S3EventRecord{
				{
					S3: events.S3Entity{
						Bucket: events.S3Bucket{
							Name: "ingestBucket",
						},
						Object: events.S3Object{
							Key: "photoKey",
						},
					},
				},
				{
					S3: events.S3Entity{
						Bucket: events.S3Bucket{
							Name: "ingestBucket",
						},
						Object: events.S3Object{
							Key: "photoKey",
						},
					},
				},
			},
		}
		s.photoRepository.On("Delete", photo.DeletePhotoParams{
			Bucket: "displayBucket",
			Key:    "photoKey",
		}).Twice().Return(nil)
		err := handler.Run(context.Background(), event)

		assert.Nil(t, err)
	})

	s.T().Run("processes s3 event and deletes photos from display bucket", func(t *testing.T) {
		s.setupMocks()
		handler := NewHandler("displayBucket", s.photoRepository)
		event := events.S3Event{
			Records: []events.S3EventRecord{
				{
					S3: events.S3Entity{
						Bucket: events.S3Bucket{
							Name: "ingestBucket",
						},
						Object: events.S3Object{
							Key: "photoKey",
						},
					},
				},
			},
		}
		s.photoRepository.On("Delete", photo.DeletePhotoParams{
			Bucket: "displayBucket",
			Key:    "photoKey",
		}).Once().Return(errors.New("something went wrong"))
		err := handler.Run(context.Background(), event)

		assert.Equal(t, "error deleting photoKey from displayBucket: something went wrong", err.Error())
	})
}

func (s *handlerTestSuite) setupMocks() {
	s.photoRepository = new(mockPhotoRepository)
}

type mockPhotoRepository struct {
	mock.Mock
}

func (m *mockPhotoRepository) Delete(params photo.DeletePhotoParams) error {
	args := m.Called(params)
	return args.Error(0)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}
