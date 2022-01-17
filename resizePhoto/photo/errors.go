package photo

type GetPhotoError struct {
	Err error
}

func (err GetPhotoError) Unwrap() error {
	return err.Err
}

func (err GetPhotoError) Error() string {
	return err.Err.Error()
}

func (err GetPhotoError) Is(target error) bool {
	_, ok := target.(GetPhotoError)
	if !ok {
		_, ok = target.(*GetPhotoError)
	}
	return ok
}

type PutPhotoError struct {
	Err error
}

func (err PutPhotoError) Unwrap() error {
	return err.Err
}

func (err PutPhotoError) Error() string {
	return err.Err.Error()
}

func (err PutPhotoError) Is(target error) bool {
	_, ok := target.(PutPhotoError)
	if !ok {
		_, ok = target.(*PutPhotoError)
	}
	return ok
}
