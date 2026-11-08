package xml

import (
	"arshin/internal/domain"
	"encoding/xml"
	"errors"
	"os"
	"strconv"
)

type converter struct{}

func NewConverter() converter {
	return converter{}
}

func (c converter) Deliver(dst string, verifications []domain.Verification) error {
	if len(verifications) == 0 {
		return errors.New("Пустой список поверок")
	}

	file, err := os.Create(dst)
	if err != nil {
		return err
	}

	application := toApplication(verifications)
	encoder := xml.NewEncoder(file)

	if err := encoder.Encode(application); err != nil {
		return err
	}

	return nil
}

func toApplication(verifications []domain.Verification) Application {
	res := Application{}

	for i := range verifications {
		record := RecInfo{
			MiInfo: MiInfo{
				SingleMI: SingleMI{
					MitypeNumber:   String32(verifications[i].Regnum()),
					ManufactureNum: String128(verifications[i].ManufactureNum()),
					Modification:   String128(verifications[i].Modification()),
				},
			},
			SignCipher:  SignCipher(verifications[i].SignCipher()),
			MiOwner:     String512(verifications[i].Owner()),
			VrfDate:     verifications[i].VerificationDate(),
			Type:        getType(verifications[i].Periodic()),
			DocTitle:    String128(verifications[i].Method()),
			Metrologist: String128(verifications[i].Verifier()),
			Means: Means{
				Mieta: Mieta{
					Number: []String32{String32(verifications[i].TestDevice().Number())},
				},
				Mis: getMis(verifications[i].TestDevice().Instruments()),
			},
			Conditions: getConditions(verifications[i].Temperature(), verifications[i].Humidity(), verifications[i].Pressure()),
		}

		if verifications[i].Fit() {
			record.ValidDate = verifications[i].ValidUntil()
		} else {
			record.Inapplicable.Reasons = "Не соответствует требованиям МП"
		}

		res.Result = append(res.Result, record)
	}

	return res
}

func getType(periodic bool) Type {
	if periodic {
		return Type(Periodic)
	}

	return Type(First)
}

func getMis(instruments []domain.MeasuringInstrument) Mis {
	var res Mis
	for _, instrument := range instruments {
		res.Mi = append(res.Mi, Mi{TypeNum: String32(instrument.Regnum()), ManufactureNum: String128(instrument.ManufactureNum())})
	}

	return res
}

func getConditions(temperature, humidity, pressure float64) Conditions {
	if temperature == 0 && humidity == 0 && pressure == 0 {
		return Conditions{
			Temperature: placeholder,
			Hymidity:    placeholder,
			Pressure:    placeholder,
		}
	}

	return Conditions{
		Temperature: String128(strconv.FormatFloat(temperature, 'f', 2, 64) + " °C"),
		Hymidity:    String128(strconv.FormatFloat(humidity, 'f', 2, 64)) + "%",
		Pressure:    String128(strconv.FormatFloat(pressure, 'f', 2, 64)) + " кПа",
	}
}
