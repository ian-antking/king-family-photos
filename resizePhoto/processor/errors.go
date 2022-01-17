package processor

type DecodeImageError struct {
	Err error
}

func (err DecodeImageError) Unwrap() error {
	return err.Err
}

func (err DecodeImageError) Error() string {
	return err.Err.Error()
}

func (err DecodeImageError) Is(target error) bool {
	_, ok := target.(DecodeImageError)
	if !ok {
		_, ok = target.(*DecodeImageError)
	}
	return ok
}

type EncodeImageError struct {
	Err error
}

func (err EncodeImageError) Unwrap() error {
	return err.Err
}

func (err EncodeImageError) Error() string {
	return err.Err.Error()
}

func (err EncodeImageError) Is(target error) bool {
	_, ok := target.(EncodeImageError)
	if !ok {
		_, ok = target.(*EncodeImageError)
	}
	return ok
}
