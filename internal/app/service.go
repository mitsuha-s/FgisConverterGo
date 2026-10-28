package app

import (
	"arshin/internal/domain"
	"log"
)

type Loader interface {
	Load(string) ([]domain.VerificationRecord, error)
}

type Sink interface {
	Send(string, []domain.Verification) error
}

type app struct {
	config *domain.Config
	logger *log.Logger
}

func New(config *domain.Config, logger *log.Logger) *app {
	return &app{config: config, logger: logger}
}

func (a *app) Handle(source string, loader Loader, dest string, sink Sink) error {
	var res []domain.Verification
	records, err := loader.Load(source)
	if err != nil {
		return err
	}

	for _, record := range records {
		if err := domain.IsValid(record); err != nil {
			a.logger.Printf("Невозможно передать счетчик %s: %v", record.ManufactureNum, err)
			continue
		}

		verification := domain.Verification{
			Fit: record.ValidUntil.IsZero(),
			Meter: domain.Meter{
				RegNum: record.RegNum,
				ManufactureNum: record.ManufactureNum,
				MIType: record.MIType,
			},
			Owner: "ФИЗИЧЕСКОЕ ЛИЦО",
			VerificationDate: record.VerificationDate,
			ValidUntil: record.ValidUntil,
			Verifier: record.Verifier,
			TestDevice: domain.TestDevice{
				Number: record.TestDeviceNum,
				Thermometer: a.config.TestDevices[record.TestDeviceNum].Thermometer,
				Stopwatch: a.config.TestDevices[record.TestDeviceNum].Stopwatch,
				Thermohygrometer: a.config.TestDevices[record.TestDeviceNum].Thermohygrometer,
			},
			Method: a.config.Meters[record.RegNum],
			SignCipher: a.config.SignCipher,
		}
		
		switch a.config.AmbientMod {
		case domain.AmbientModFile: 
			verification.Temperature = record.Temperature
			verification.Humidity = record.Humidity
			verification.Pressure = record.Pressure
		case domain.AmbientModRandom:
			temp, humidity, pressure := getRandomAmbient()
			verification.Temperature = temp
			verification.Humidity = humidity
			verification.Pressure = pressure
		}

		res = append(res, verification)
	}

	if err := sink.Send(dest, res); err != nil {
		return err
	}

	return nil
}

func getRandomAmbient() (temp float64, humidity float64, pressure float64) {
	// TODO
	return
}