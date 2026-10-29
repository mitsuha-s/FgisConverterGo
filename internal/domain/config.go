package domain

type SourceAmbient int

const (
	AmbientModEmpty SourceAmbient = iota
	AmbientModFile
	AmbientModRandom
)

type Config struct {
	SignCipher  string
	Meters      map[string]string
	TestDevices map[string]TestDevice
	AmbientMod  SourceAmbient
}
