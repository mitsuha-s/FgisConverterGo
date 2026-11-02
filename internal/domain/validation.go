package domain

import (
	"regexp"
)

func ValidateRegnum(regnum string) error {
	ok, err := regexp.MatchString(patternRegNum, regnum)
	if err != nil {
		return err
	}
	if !ok {
		return ErrInvalidData{"Некорректный регистрационный номер"}
	}

	return nil
}
