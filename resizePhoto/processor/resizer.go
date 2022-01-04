package processor

type Resizer struct {
	width  uint
	height uint
}

func (r *Resizer) Run(image Image) (Image, error) {
	return image, nil
}

func NewResizer(width, height uint) Resizer {
	return Resizer{
		width:  width,
		height: height,
	}
}
