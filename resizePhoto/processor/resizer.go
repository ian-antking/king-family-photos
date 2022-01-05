package processor

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"

	"github.com/nfnt/resize"
)

type Resizer struct {
	width  uint
	height uint
}

func (r *Resizer) Run(imageInput Image) (Image, error) {
	img, _, decodeErr := image.Decode(bytes.NewReader(imageInput.Image))
	if nil != decodeErr {
		return Image{}, fmt.Errorf("error decoding image: %s", decodeErr.Error())
	}

	resizedImage := resize.Resize(r.width, r.height, img, resize.Lanczos3)

	buffer := new(bytes.Buffer)
	encodeErr := jpeg.Encode(buffer, resizedImage, nil)

	if nil != encodeErr {
		return Image{}, fmt.Errorf("error encoding image: %s", decodeErr.Error())
	}

	newImage := Image{
		Image:  buffer.Bytes(),
		Bucket: imageInput.Bucket,
		Key:    imageInput.Key,
	}

	return newImage, nil
}

func NewResizer(width, height uint) Resizer {
	return Resizer{
		width:  width,
		height: height,
	}
}
