package domain

import "time"

type Verification struct {
	Fit              bool
	Meter            Meter
	Owner            string
	VerificationDate time.Time
	ValidUntil       time.Time
	Verifier         string
	TestDevice       TestDevice
	Method           string
	Temperature      float64
	Humidity         float64
	Pressure         float64
	SignCipher       string
}

type Meter struct {
	RegNum         string
	ManufactureNum string
	MIType         string
}

type TestDevice struct {
	Number           string
	Stopwatch        MeasuringInstrument
	Thermometer      MeasuringInstrument
	Thermohygrometer MeasuringInstrument
}

type MeasuringInstrument struct {
	RegNum         string
	ManufactureNum string
}
