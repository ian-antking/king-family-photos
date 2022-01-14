package integration

import (
	"image"
	"time"

	"github.com/stretchr/testify/assert"
)

func (s *integrationTestSuite) TestUploadImages() {
	img := image.NewRGBA(image.Rect(0, 0, 3648, 2736))

	objectKey := s.putImageInIngestBucket(img)

	time.Sleep(time.Second * 5)

	objects := s.listItemsInBucket(s.displayBucketName)

	assert.Contains(s.T(), objects, objectKey)
}
