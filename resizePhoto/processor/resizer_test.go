package processor

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type resizerTestSuite struct {
	suite.Suite
}

func (s *resizerTestSuite) TestRun() {
	s.T().Run("reduces image width/height to values set at instantiation", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 100, 100))

		buf := new(bytes.Buffer)
		_ = jpeg.Encode(buf, img, nil)

		resizer := NewResizer(50, 50)

		output, err := resizer.Run(Image{
			Image:  buf.Bytes(),
			Bucket: "bucket",
			Key:    "key",
		})

		assert.Nil(t, err)

		resizedImage, _, _ := image.Decode(bytes.NewReader(output.Image))

		assert.Equal(t, 50, resizedImage.Bounds().Max.X)
		assert.Equal(t, 50, resizedImage.Bounds().Max.Y)
	})

	s.T().Run("increases image width/height to values set at instantiation", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 100, 100))

		buf := new(bytes.Buffer)
		_ = jpeg.Encode(buf, img, nil)

		resizer := NewResizer(200, 200)

		output, err := resizer.Run(Image{
			Image:  buf.Bytes(),
			Bucket: "bucket",
			Key:    "key",
		})

		assert.Nil(t, err)

		resizedImage, _, _ := image.Decode(bytes.NewReader(output.Image))

		assert.Equal(t, 200, resizedImage.Bounds().Max.X)
		assert.Equal(t, 200, resizedImage.Bounds().Max.Y)
	})

	s.T().Run("decreasing size maintains width/height ratio if only one parameter set", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 200, 100))

		buf := new(bytes.Buffer)
		_ = jpeg.Encode(buf, img, nil)

		resizer := NewResizer(0, 50)

		output, err := resizer.Run(Image{
			Image:  buf.Bytes(),
			Bucket: "bucket",
			Key:    "key",
		})

		assert.Nil(t, err)

		resizedImage, _, _ := image.Decode(bytes.NewReader(output.Image))

		assert.Equal(t, 100, resizedImage.Bounds().Max.X)
		assert.Equal(t, 50, resizedImage.Bounds().Max.Y)
	})

	s.T().Run("increasing size maintains width/height ratio if only one parameter set", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 200, 100))

		buf := new(bytes.Buffer)
		_ = jpeg.Encode(buf, img, nil)

		resizer := NewResizer(0, 100)

		output, err := resizer.Run(Image{
			Image:  buf.Bytes(),
			Bucket: "bucket",
			Key:    "key",
		})

		assert.Nil(t, err)

		resizedImage, _, _ := image.Decode(bytes.NewReader(output.Image))

		assert.Equal(t, 200, resizedImage.Bounds().Max.X)
		assert.Equal(t, 100, resizedImage.Bounds().Max.Y)
	})

	s.T().Run("handles png encoded images", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 100, 100))

		buf := new(bytes.Buffer)
		_ = png.Encode(buf, img)

		resizer := NewResizer(50, 50)

		output, err := resizer.Run(Image{
			Image:  buf.Bytes(),
			Bucket: "bucket",
			Key:    "key",
		})

		assert.Nil(t, err)

		resizedImage, _, _ := image.Decode(bytes.NewReader(output.Image))

		assert.Equal(t, 50, resizedImage.Bounds().Max.X)
		assert.Equal(t, 50, resizedImage.Bounds().Max.Y)
	})
}

func TestResizerTestSuite(t *testing.T) {
	suite.Run(t, new(resizerTestSuite))
}
