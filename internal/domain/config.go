package domain

import (
	"maps"
	"strings"
)

type SourceAmbient int

const (
	AmbientModeEmpty SourceAmbient = iota + 1
	AmbientModeFile
	AmbientModeRandom
)

type Config struct {
	signCipher  string
	methods     map[string]bool
	mitypes     map[regnum]string
	testDevices map[string]TestDevice
	ambientMode SourceAmbient
}

func NewConfig(signCipher string, mitypes map[regnum]string, methods []string, testDevices []TestDevice, ambientMode SourceAmbient) (*Config, error) {
	if signCipher == "" {
		return nil, ErrInvalidConfig{"Не указан шифр организации"}
	}

	if ambientMode == 0 {
		return nil, ErrInvalidConfig{"Не указан режим обработки показаний окружающей среды"}
	}

	devicesmap, err := getTestDevicesMap(testDevices)
	if err != nil {
		return nil, err
	}

	methodsmap, err := getMethodsMap(methods)
	if err != nil {
		return nil, err
	}

	if err := validateMitypes(mitypes, methodsmap); err != nil {
		return nil, err
	}

	mitypesCopy := make(map[regnum]string)
	maps.Copy(mitypesCopy, mitypes)

	return &Config{
		signCipher:  signCipher,
		methods:     methodsmap,
		mitypes:     mitypesCopy,
		testDevices: devicesmap,
		ambientMode: ambientMode,
	}, nil
}

func getTestDevicesMap(devices []TestDevice) (map[string]TestDevice, error) {
	res := make(map[string]TestDevice)

	if len(devices) == 0 {
		return nil, ErrInvalidConfig{"Отсутствуют поверочные установки"}
	}

	for _, device := range devices {
		if device == (TestDevice{}) {
			return nil, ErrInvalidConfig{"Пустая структура поверочной установки"}
		}
		res[device.number] = device
	}

	return res, nil
}

func getMethodsMap(methods []string) (map[string]bool, error) {
	res := make(map[string]bool)

	if len(methods) == 0 {
		return nil, ErrInvalidConfig{"Отсустствуют разрешенные методики поверки"}
	}

	for _, method := range methods {
		if method == "" {
			return nil, ErrInvalidConfig{"Методика поверки не может быть пустой строкой"}
		}
		res[method] = true
	}

	return res, nil
}

func validateMitypes(mitypes map[regnum]string, methods map[string]bool) error {
	if len(mitypes) == 0 {
		return ErrInvalidConfig{"Отсутствует список поверяемых типов СИ"}
	}

	for num, method := range mitypes {
		if !methods[method] {
			errstring := getMitypeErrorsting(num, method, methods)
			return ErrInvalidConfig{errstring}
		}
	}

	return nil
}

func getMitypeErrorsting(num regnum, method string, methods map[string]bool) string {
	var b strings.Builder

	b.WriteString("Несуществующая методика ")
	b.WriteString(method)
	b.WriteString(" у счетчика ")
	b.WriteString(string(num))
	b.WriteString(". Список разрешенных методик:\n")

	for method := range methods {
		b.WriteString(method)
		b.WriteByte('\n')
	}

	return b.String()
}
