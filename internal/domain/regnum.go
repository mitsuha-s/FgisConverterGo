package domain

import (
	"errors"
	"regexp"
)

type regnum string

func NewRegnum(num string) (regnum, error) {
	ok, err := regexp.Match(`^\d{5}-\d{2}$`, []byte(num))
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("Некорректный регистрационный номер")
	}

	return regnum(num), nil
}
