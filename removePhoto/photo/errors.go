package photo

type DeletePhotoError struct {
	Err error
}

func (err DeletePhotoError) Unwrap() error {
	return err.Err
}

func (err DeletePhotoError) Error() string {
	return err.Err.Error()
}

func (err DeletePhotoError) Is(target error) bool {
	_, ok := target.(DeletePhotoError)
	if !ok {
		_, ok = target.(*DeletePhotoError)
	}
	return ok
}
