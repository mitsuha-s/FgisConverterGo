package domain

type ErrInvalidData struct {
	message string
}

func (e ErrInvalidData) Error() string {
	return e.message
}
