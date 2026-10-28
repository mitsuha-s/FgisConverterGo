package domain

import "time"

type VerificationRecord struct {
	RegNum           string
	ManufactureNum   string
	MIType           string
	VerificationDate time.Time
	ValidUntil       time.Time
	Verifier         string
	TestDeviceNum    string
	Temperature      string
	Humidity         string
	Pressure         string
}
