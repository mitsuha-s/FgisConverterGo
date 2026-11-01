package domain

import (
	"time"
)

type VerificationRecord struct {
	RegNum           string
	ManufactureNum   string
	MIType           string
	VerificationDate time.Time
	ValidUntil       time.Time
	Verifier         string
	TestDeviceNum    string
	Temperature      float64
	Humidity         float64
	Pressure         float64
}
