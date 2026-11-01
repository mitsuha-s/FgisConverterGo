package domain

type ErrInvalidConfig struct {
	message string
}

func (e ErrInvalidConfig) Error() string {
	return e.message
}

type ErrInvalidData struct {
	message string
}

func (e ErrInvalidData) Error() string {
	return e.message
}
