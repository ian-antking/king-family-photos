package processor

type Image struct {
	Image  []byte
	Bucket string
	Key    string
}

type Processor interface {
	Run(Image) (Image, error)
}
