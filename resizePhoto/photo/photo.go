package photo

type GetPhotoParams struct {
	Bucket string
	Key    string
}
type GetPhotoOutput struct {
	Image  []byte
	Bucket string
	Key    string
}

type PutPhotoParams struct {
	Image  []byte
	Key    string
	Bucket string
}

type Repository interface {
	Get(GetPhotoParams) (GetPhotoOutput, error)
	Put(PutPhotoParams) error
}
