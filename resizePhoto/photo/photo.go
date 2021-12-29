package photo

import "io"

type GetPhotoParams struct {
	Bucket string
	Key string
}
type GetPhotoOutput struct {
	Image io.ReadCloser
}

type Repository interface {
	Get(GetPhotoParams) (GetPhotoOutput, error)
}
