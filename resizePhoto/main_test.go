package main

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ian-antking/king-family-photos/resizePhoto/photo"
	"github.com/ian-antking/king-family-photos/resizePhoto/processor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type handlerTestSuite struct {
	suite.Suite
	photoRepository *mockPhotoRepository
	imageProcessor  *mockImageProcessor
}

func (s *handlerTestSuite) setUpMocks() {
	s.photoRepository = new(mockPhotoRepository)
	s.imageProcessor = new(mockImageProcessor)
}

func (s *handlerTestSuite) TestGetPhotoParams() {
	s.T().Run("extracts photo bucket names and keys from s3 event records", func(t *testing.T) {
		event := events.S3Event{
			Records: []events.S3EventRecord{
				{
					S3: events.S3Entity{
						Bucket: events.S3Bucket{
							Name: "bucketName",
						},
						Object: events.S3Object{
							Key: "photo1",
						},
					},
				},
				{
					S3: events.S3Entity{
						Bucket: events.S3Bucket{
							Name: "bucketName",
						},
						Object: events.S3Object{
							Key: "photo2",
						},
					},
				},
			},
		}
		expected := []photo.GetPhotoParams{
			{
				Bucket: "bucketName",
				Key:    "photo1",
			},
			{
				Bucket: "bucketName",
				Key:    "photo2",
			},
		}

		actual := getPhotoParams(event)

		assert.Equal(t, expected, actual)
	})
}

func (s *handlerTestSuite) TestGetImages() {
	s.T().Run("returns slice of images from s3", func(t *testing.T) {
		s.setUpMocks()
		handler := NewHandler(s.photoRepository, "bucket", s.imageProcessor)
		s.photoRepository.On("Get", photo.GetPhotoParams{
			Bucket: "bucket",
			Key:    "photo",
		}).Return(photo.GetPhotoOutput{
			Image:  []byte{},
			Bucket: "bucket",
			Key:    "photo",
		}, nil)
		expected := []photo.GetPhotoOutput{
			{
				Image:  []byte{},
				Bucket: "bucket",
				Key:    "photo",
			},
		}

		actual, err := handler.getImages([]photo.GetPhotoParams{
			{
				Bucket: "bucket",
				Key:    "photo",
			},
		})

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	s.T().Run("returns error if failed to get image from s3", func(t *testing.T) {
		s.setUpMocks()
		handler := NewHandler(s.photoRepository, "bucket", s.imageProcessor)

		s.photoRepository.On("Get", photo.GetPhotoParams{
			Bucket: "bucket",
			Key:    "photo1",
		}).Return(photo.GetPhotoOutput{
			Image:  []byte{},
			Bucket: "bucket",
			Key:    "photo1",
		}, nil)
		s.photoRepository.On("Get", photo.GetPhotoParams{
			Bucket: "bucket",
			Key:    "photo2",
		}).Return(photo.GetPhotoOutput{}, errors.New("something went wrong"))

		_, err := handler.getImages([]photo.GetPhotoParams{
			{
				Bucket: "bucket",
				Key:    "photo1",
			},
			{
				Bucket: "bucket",
				Key:    "photo2",
			},
		})

		assert.Equal(t, "something went wrong", err.Error())
	})
}

func (s *handlerTestSuite) TestProcessImages() {
	s.T().Run("returns slice of processed images", func(t *testing.T) {
		s.setUpMocks()
		handler := NewHandler(s.photoRepository, "bucket", s.imageProcessor)
		expected := []processor.Image{
			{
				Image:  []byte{},
				Bucket: "bucket",
				Key:    "photo",
			},
			{
				Image:  []byte{},
				Bucket: "bucket",
				Key:    "photo",
			},
		}

		s.imageProcessor.On("Run", processor.Image{
			Image:  []byte{},
			Bucket: "bucket",
			Key:    "photo",
		}).Return(processor.Image{
			Image:  []byte{},
			Bucket: "bucket",
			Key:    "photo",
		}, nil)

		actual, err := handler.processImages([]photo.GetPhotoOutput{
			{
				Image:  []byte{},
				Bucket: "bucket",
				Key:    "photo",
			}, {
				Image:  []byte{},
				Bucket: "bucket",
				Key:    "photo",
			},
		})

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	s.T().Run("return error if image failed to process", func(t *testing.T) {
		s.setUpMocks()
		handler := NewHandler(s.photoRepository, "bucket", s.imageProcessor)

		s.imageProcessor.On("Run", processor.Image{
			Image:  []byte{},
			Bucket: "bucket",
			Key:    "photo1",
		}).Return(processor.Image{
			Image:  []byte{},
			Bucket: "bucket",
			Key:    "photo1",
		}, nil)

		s.imageProcessor.On("Run", processor.Image{
			Image:  []byte{},
			Bucket: "bucket",
			Key:    "photo2",
		}).Return(processor.Image{}, errors.New("something went wrong"))

		_, err := handler.processImages([]photo.GetPhotoOutput{
			{
				Image:  []byte{},
				Bucket: "bucket",
				Key:    "photo1",
			}, {
				Image:  []byte{},
				Bucket: "bucket",
				Key:    "photo2",
			},
		})

		assert.Equal(t, "something went wrong", err.Error())
	})
}

func (s *handlerTestSuite) TestPutImages() {
	s.T().Run("returns error if image failed to upload", func(t *testing.T) {
		s.setUpMocks()
		handler := NewHandler(s.photoRepository, "bucket", s.imageProcessor)

		s.photoRepository.On("Put", photo.PutPhotoParams{
			Image:  []byte{},
			Bucket: "bucket",
			Key:    "photo1",
		}).Return(nil)
		s.photoRepository.On("Put", photo.PutPhotoParams{
			Image:  []byte{},
			Bucket: "bucket",
			Key:    "photo2",
		}).Return(errors.New("something went wrong"))

		err := handler.putImages([]processor.Image{
			{
				Image:  []byte{},
				Bucket: "bucket",
				Key:    "photo1",
			},
			{
				Image:  []byte{},
				Bucket: "bucket",
				Key:    "photo2",
			},
		})

		assert.Equal(t, "something went wrong", err.Error())
	})
}

func (s *handlerTestSuite) TestRun() {
	s.T().Run("returns s3.Get error", func(t *testing.T) {
		s.setUpMocks()
		handler := NewHandler(s.photoRepository, "bucket", s.imageProcessor)

		s.photoRepository.On("Get", mock.Anything).Return(photo.GetPhotoOutput{}, errors.New("something went wrong"))

		err := handler.Run(context.Background(), events.S3Event{
			Records: []events.S3EventRecord{
				{
					S3: events.S3Entity{
						Bucket: events.S3Bucket{
							Name: "bucket",
						},
						Object: events.S3Object{
							Key: "photo",
						},
					},
				},
			},
		})

		assert.NotNil(t, err)
		assert.Equal(t, "something went wrong", err.Error())
	})

	s.T().Run("returns processor.Run error", func(t *testing.T) {
		s.setUpMocks()
		handler := NewHandler(s.photoRepository, "bucket", s.imageProcessor)

		s.photoRepository.On("Get", mock.Anything).Return(photo.GetPhotoOutput{
			Image:  []byte{},
			Bucket: "bucket",
			Key:    "photo",
		}, nil)

		s.imageProcessor.On("Run", mock.Anything).Return(processor.Image{}, errors.New("something went wrong"))

		err := handler.Run(context.Background(), events.S3Event{
			Records: []events.S3EventRecord{
				{
					S3: events.S3Entity{
						Bucket: events.S3Bucket{
							Name: "bucket",
						},
						Object: events.S3Object{
							Key: "photo",
						},
					},
				},
			},
		})

		assert.NotNil(t, err)
		assert.Equal(t, "something went wrong", err.Error())
	})

	s.T().Run("returns s3.Put error", func(t *testing.T) {
		s.setUpMocks()
		handler := NewHandler(s.photoRepository, "bucket", s.imageProcessor)

		s.photoRepository.On("Get", mock.Anything).Return(photo.GetPhotoOutput{
			Image:  []byte{},
			Bucket: "bucket",
			Key:    "photo",
		}, nil)

		s.imageProcessor.On("Run", mock.Anything).Return(processor.Image{
			Image:  []byte{},
			Bucket: "bucket",
			Key:    "photo",
		}, nil)

		s.photoRepository.On("Put", mock.Anything).Return(errors.New("something went wrong"))

		err := handler.Run(context.Background(), events.S3Event{
			Records: []events.S3EventRecord{
				{
					S3: events.S3Entity{
						Bucket: events.S3Bucket{
							Name: "bucket",
						},
						Object: events.S3Object{
							Key: "photo",
						},
					},
				},
			},
		})

		assert.NotNil(t, err)
		assert.Equal(t, "something went wrong", err.Error())
	})
}

func (s *handlerTestSuite) setupMocks() {
	s.photoRepository = new(mockPhotoRepository)
}

type mockPhotoRepository struct {
	mock.Mock
}

func (m *mockPhotoRepository) Get(params photo.GetPhotoParams) (photo.GetPhotoOutput, error) {
	args := m.Called(params)
	return args.Get(0).(photo.GetPhotoOutput), args.Error(1)
}

func (m *mockPhotoRepository) Put(params photo.PutPhotoParams) error {
	args := m.Called(params)
	return args.Error(0)
}

type mockImageProcessor struct {
	mock.Mock
}

func (m *mockImageProcessor) Run(image processor.Image) (processor.Image, error) {
	args := m.Called(image)
	return args.Get(0).(processor.Image), args.Error(1)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}
