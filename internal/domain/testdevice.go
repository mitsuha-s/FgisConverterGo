package domain

type TestDevice struct {
	number      string
	instruments []MeasuringInstrument
}

func NewTestDevice(number string, instruments []MeasuringInstrument) (TestDevice, error) {
	if number == "" {
		return TestDevice{}, ErrInvalidData{"Не указан заводской номер"}
	}

	instrumentsCopy := make([]MeasuringInstrument, len(instruments))
	copy(instrumentsCopy, instruments)

	if err := validateInstruments(instrumentsCopy); err != nil {
		return TestDevice{}, err
	}

	return TestDevice{number: number, instruments: instrumentsCopy}, nil
}

func (t TestDevice) IsZero() bool {
	if t.number == "" && t.instruments == nil {
		return true
	}

	return false
}

func validateInstruments(instruments []MeasuringInstrument) error {
	if len(instruments) == 0 {
		return ErrInvalidData{"Не указаны доп.средства измерения"}
	}

	for _, instrument := range instruments {
		if instrument == (MeasuringInstrument{}) {
			return ErrInvalidData{"Пустая структура доп.средства измерения"}
		}
	}

	return nil
}
