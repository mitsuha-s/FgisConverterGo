package domain

type MeasuringInstrument struct {
	regNum         regnum
	manufactureNum string
}

func NewMeasuringInstrument(reg regnum, manufactureNum string) (MeasuringInstrument, error) {
	if manufactureNum == "" {
		return MeasuringInstrument{}, ErrInvalidData{"Не указан заводской номер дополнительного СИ"}
	}

	return MeasuringInstrument{regNum: reg, manufactureNum: manufactureNum}, nil
}
