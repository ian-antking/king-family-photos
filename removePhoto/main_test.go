package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/ian-antking/king-family-photos/removePhoto/photo"
)

type handlerTestSuite struct {
	suite.Suite
	photoRepository photo.Repository
}

func (s *handlerTestSuite) TestGetPhotoParams() {
	s.T().Run("converts records on an S3Event to DeletePhotoParams", func(t *testing.T) {
		s.setupMocks()
		handler := NewHandler("displayBucket", s.photoRepository)
		expected := []photo.DeletePhotoParams{
			{
				Bucket: "displayBucket",
				Key:    "photoKey",
			},
		}

		result := handler.getPhotoParams(events.S3Event{
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
		})

		assert.Equal(t, expected, result)
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
