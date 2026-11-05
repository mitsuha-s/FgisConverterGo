package domain

import (
	"fmt"
	"time"
)

type Verification struct {
	fit              bool
	periodic         bool
	regNum           string
	manufactureNum   string
	modification     string
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

	err := ValidateRegnum(record.RegNum)
	if err != nil {
		return Verification{}, err
	}

	method, ok := config.mitypes[record.RegNum]
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
		periodic:         true,
		regNum:           record.RegNum,
		manufactureNum:   record.ManufactureNum,
		modification:     record.Modification,
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

func (v *Verification) Fit() bool {
	return v.fit
}

func (v *Verification) Periodic() bool {
	return v.periodic
}

func (v *Verification) Regnum() string {
	return v.regNum
}

func (v *Verification) ManufactureNum() string {
	return v.manufactureNum
}

func (v *Verification) Modification() string {
	return v.modification
}

func (v *Verification) Owner() string {
	return v.owner
}

func (v *Verification) VerificationDate() time.Time {
	return v.verificationDate
}

func (v *Verification) ValidUntil() time.Time {
	return v.validUntil
}

func (v *Verification) Verifier() string {
	return v.verifier
}

func (v *Verification) TestDevice() TestDevice {
	return v.testDevice
}

func (v *Verification) Method() string {
	return v.method
}

func (v *Verification) Temperature() float64 {
	return v.temperature
}

func (v *Verification) Humidity() float64 {
	return v.humidity
}

func (v *Verification) Pressure() float64 {
	return v.pressure
}

func (v *Verification) SignCipher() string {
	return v.signCipher
}

func NewTestVirifications() []Verification {
	res := []Verification{
		Verification{
			fit:              true,
			periodic:         true,
			regNum:           "12345-67",
			manufactureNum:   "1233",
			modification:     "TEST-15",
			owner:            "ТЕСТОВОЕ ЛИЦО",
			verificationDate: time.Now(),
			validUntil:       time.Date(2026, time.December, 31, 0, 0, 0, 0, time.Local),
			verifier:         "Тестов Т.Т.",
			testDevice: TestDevice{
				number: "1234.3P.456",
				instruments: []MeasuringInstrument{
					MeasuringInstrument{
						regNum:         "34568-90",
						manufactureNum: "1234654",
					},
				},
			},
			method:     "МИ 2015-2019",
			signCipher: "ДГХ",
		},
		Verification{
			fit:              false,
			periodic:         true,
			regNum:           "12345-67",
			manufactureNum:   "НЕ ГОДЕН",
			modification:     "TEST-15",
			owner:            "ТЕСТОВОЕ ЛИЦО",
			verificationDate: time.Now(),
			validUntil:       time.Date(2026, time.December, 31, 0, 0, 0, 0, time.Local),
			verifier:         "Тестов Т.Т.",
			testDevice: TestDevice{
				number: "1234.3P.456",
				instruments: []MeasuringInstrument{
					MeasuringInstrument{
						regNum:         "34568-90",
						manufactureNum: "1234654",
					},
				},
			},
			method:      "МИ 2015-2019",
			signCipher:  "ДГХ",
			temperature: 35,
			humidity:    46,
			pressure:    100.2,
		},
	}

	return res
}
