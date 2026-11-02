package domain

type MeasuringInstrument struct {
	regNum         string
	manufactureNum string
}

func NewMeasuringInstrument(reg string, manufactureNum string) (MeasuringInstrument, error) {
	if manufactureNum == "" {
		return MeasuringInstrument{}, ErrInvalidData{"Не указан заводской номер"}
	}
	if err := ValidateRegnum(reg); err != nil {
		return MeasuringInstrument{}, err
	}

	return MeasuringInstrument{regNum: reg, manufactureNum: manufactureNum}, nil
}
