package domain

import (
	"fmt"
	"time"
)

type Verification struct {
	fit              bool
	regNum           string
	manufactureNum   string
	mitype           string
	owner            string
	verificationDate time.Time
	validUntil       time.Time
	verifier         string
	testDevice       TestDevice
	method           string
	// string?
	temperature float64
	humidity    float64
	pressure    float64
	signCipher  string
}

func NewVerification(record VerificationRecord, config *Config) (Verification, error) {
	var fit bool
	var temp, humidity, pressure float64

	reg, err := NewRegnum(record.RegNum)
	if err != nil {
		return Verification{}, err
	}

	method, ok := config.mitypes[reg]
	if !ok {
		return Verification{}, fmt.Errorf("Отсутствует в списке поверяемых типов СИ")
	}

	if !record.ValidUntil.IsZero() {
		if record.VerificationDate.Compare(record.ValidUntil) != -1 {
			return Verification{}, fmt.Errorf("Дата поверки должна предшествовать дате последнего валидного дня")
		}

		verificationInterval := record.ValidUntil.Year() - record.VerificationDate.Year()

		if verificationInterval < MinVerificationInterval || verificationInterval > MaxVerificationIntrval {
			return Verification{}, fmt.Errorf("Некорректный межповерочный интервал: %d лет. МПИ может быть только от 4 по 6 лет включительно", verificationInterval)
		}

		fit = true
	}

	testDevice, ok := config.testDevices[record.TestDeviceNum]
	if !ok {
		return Verification{}, fmt.Errorf("Установка %s не найдена", record.TestDeviceNum)
	}

	switch config.ambientMode {
	case AmbientModeFile:
		if record.Humidity == 0 {
			return Verification{}, fmt.Errorf("Влажность не указана или равна нулю")
		}
		if record.Temperature == 0 {
			return Verification{}, fmt.Errorf("Температура не указана или равна нулю")
		}
		if record.Pressure == 0 {
			return Verification{}, fmt.Errorf("Атм.давление не указано или равно нулю")
		}

		temp = record.Temperature
		humidity = record.Humidity
		pressure = record.Pressure
	case AmbientModeRandom:
		temp, humidity, pressure = getRandomAmbient()
	}

	return Verification{
		fit:              fit,
		regNum:           record.RegNum,
		manufactureNum:   record.ManufactureNum,
		mitype:           record.MIType,
		owner:            "ФИЗИЧЕСКОЕ ЛИЦО",
		verificationDate: record.VerificationDate,
		validUntil:       record.ValidUntil,
		verifier:         record.Verifier,
		testDevice:       testDevice,
		method:           method,
		temperature:      temp,
		humidity:         humidity,
		pressure:         pressure,
		signCipher:       config.signCipher,
	}, nil
}

func getRandomAmbient() (temp float64, humidity float64, pressure float64) {
	// TODO
	return
}
