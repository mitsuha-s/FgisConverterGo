package domain

import (
	"fmt"
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
	mitypes     map[string]string
	testDevices map[string]TestDevice
	ambientMode SourceAmbient
}

func NewConfig(signCipher string, mitypes map[string]string, methods []string, testDevices []TestDevice, ambientMode SourceAmbient) (*Config, error) {
	if signCipher == "" {
		return nil, ErrInvalidData{"Не указан шифр организации"}
	}

	if ambientMode == 0 {
		return nil, ErrInvalidData{"Не указан режим обработки показаний окружающей среды"}
	}

	devicesmap, err := getTestDevicesMap(testDevices)
	if err != nil {
		return nil, err
	}

	methodsmap, err := getMethodsMap(methods)
	if err != nil {
		return nil, err
	}

	mitypesCopy := make(map[string]string)
	maps.Copy(mitypesCopy, mitypes)

	if err := validateMitypes(mitypesCopy, methodsmap); err != nil {
		return nil, err
	}

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
		return nil, ErrInvalidData{"Отсутствуют поверочные установки"}
	}

	for _, device := range devices {
		if device.IsZero() {
			return nil, ErrInvalidData{"Пустая структура поверочной установки"}
		}
		res[device.number] = device
	}

	return res, nil
}

func getMethodsMap(methods []string) (map[string]bool, error) {
	res := make(map[string]bool)

	if len(methods) == 0 {
		return nil, ErrInvalidData{"Отсутствуют разрешенные методики поверки"}
	}

	for _, method := range methods {
		if method == "" {
			return nil, ErrInvalidData{"Методика поверки не может быть пустой строкой"}
		}
		res[method] = true
	}

	return res, nil
}

func validateMitypes(mitypes map[string]string, methods map[string]bool) error {
	if len(mitypes) == 0 {
		return ErrInvalidData{"Отсутствует список поверяемых типов СИ"}
	}

	for regnum, method := range mitypes {
		if err := ValidateRegnum(regnum); err != nil {
			return fmt.Errorf("Некорректный поверяемый тип СИ %s: %w", regnum, err)
		}

		if !methods[method] {
			errstring := getMitypeErrorstring(regnum, method, methods)
			return ErrInvalidData{errstring}
		}
	}

	return nil
}

func getMitypeErrorstring(regnum string, method string, methods map[string]bool) string {
	var b strings.Builder

	b.WriteString("Несуществующая методика ")
	b.WriteString(method)
	b.WriteString(" у счетчика ")
	b.WriteString(regnum)
	b.WriteString(". Список разрешенных методик:\n")

	for method := range methods {
		b.WriteString(method)
		b.WriteByte('\n')
	}

	return b.String()
}
