package domain

import (
	"errors"
	"fmt"
)

const (
	MinVerificationInterval = 4
	MaxVerificationIntrval  = 6
)

func IsValid(record VerificationRecord, config *Config) error {
	if _, ok := config.Mitypes[record.RegNum]; !ok {
		return errors.New("Невалидная методика поверки")
	}

	if !record.ValidUntil.IsZero() {
		if record.VerificationDate.Compare(record.ValidUntil) == -1 {
			return fmt.Errorf("Дата поверки должна предшествовать дате следующей поверки")
		}

		verificationInterval := record.ValidUntil.Year() - record.VerificationDate.Year()

		if verificationInterval < MinVerificationInterval || verificationInterval > MaxVerificationIntrval {
			return fmt.Errorf("Некорректный межповерочный интервал: %d лет. МПИ может быть только от 4 по 6 лет включительно", verificationInterval)
		}
	}

	if _, ok := config.TestDevices[record.ManufactureNum]; !ok {
		return fmt.Errorf("Установка %s не найдена", record.TestDeviceNum)
	}

	if config.AmbientMode == AmbientModeFile {
		if record.Humidity == 0 {
			return errors.New("Влажность не указана или равна нулю")
		}
		if record.Temperature == 0 {
			return errors.New("Температура не указана или равна нулю")
		}
		if record.Pressure == 0 {
			return errors.New("Атм. давление не указано или равно нулю")
		}
	}

	return nil
}
