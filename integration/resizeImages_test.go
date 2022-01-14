package integration

import (
	"image"
	"time"

	"github.com/stretchr/testify/assert"
)

func (s *integrationTestSuite) TestResizeImages() {
	img := image.NewRGBA(image.Rect(0, 0, 3648, 2736))

	objectKey := s.putImageInIngestBucket(img)

	time.Sleep(time.Second * 5)

	resizedImage := s.getImageFromBucket(s.displayBucketName, objectKey)

	assert.Equal(s.T(), 480, resizedImage.Bounds().Max.Y)
}
