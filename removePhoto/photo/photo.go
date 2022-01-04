package photo

type DeletePhotoParams struct {
	Bucket string
	Key    string
}

type Repository interface {
	Delete(DeletePhotoParams) error
}
