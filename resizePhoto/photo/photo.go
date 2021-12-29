package photo

type GetPhotoParams struct {
	Bucket string
	Key string
}
type GetPhotoOutput struct {}

type Repository interface {
	Get(GetPhotoParams) (GetPhotoOutput, error)
}
