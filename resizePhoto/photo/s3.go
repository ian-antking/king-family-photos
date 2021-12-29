package photo

type S3 struct {}

func (s *S3) Get(params GetPhotoParams) (GetPhotoOutput, error) {
	return GetPhotoOutput{}, nil
}
