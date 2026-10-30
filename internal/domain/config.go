package domain

type SourceAmbient int

const (
	AmbientModeEmpty SourceAmbient = iota
	AmbientModeFile
	AmbientModeRandom
)

type Config struct {
	SignCipher  string
	Mitypes     map[string]string
	TestDevices map[string]TestDevice
	AmbientMode SourceAmbient
}
