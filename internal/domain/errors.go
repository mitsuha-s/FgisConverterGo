package domain

import "errors"

type ErrInvalidData struct {
	message string
}

func (e ErrInvalidData) Error() string {
	return e.message
}

var ErrNoRecords = errors.New("Не найдено ни одной записи о поверке")
