package domain

type TestDevice struct {
	number           string
	stopwatch        MeasuringInstrument
	thermometer      MeasuringInstrument
	thermohygrometer MeasuringInstrument
}

func NewTestDevice(number string, stopwatch, thermometer, thermohygrometer MeasuringInstrument) (TestDevice, error) {
	if number == "" {
		return TestDevice{}, ErrInvalidData{"Не указан номер поверочной установки"}
	}

	if stopwatch == (MeasuringInstrument{}) && thermometer == (MeasuringInstrument{}) && thermohygrometer == (MeasuringInstrument{}) {
		return TestDevice{}, ErrInvalidData{"Не указаны доп.средства измерения для установки"}
	}

	return TestDevice{number: number, stopwatch: stopwatch, thermometer: thermometer, thermohygrometer: thermohygrometer}, nil
}
