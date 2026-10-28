package xls

import (
	"arshin/internal/domain"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

var dataLayouts = []string{"02.01.2006", "02/01/2006"}

type parser struct {
	logger *log.Logger
	opts   ParseOptions
}

func New(logger *log.Logger, opts ...ParseOptions) *parser {
	res := &parser{logger: logger}

	if len(opts) > 0 {
		res.opts = opts[0]
	}

	return res
}

func (p *parser) Load(source string) ([]domain.VerificationRecord, error) {
	var res []domain.VerificationRecord

	file, err := excelize.OpenFile(source)
	if err != nil {
		return nil, err
	}

	for _, sheet := range file.GetSheetList() {
		rows, _ := file.GetRows(sheet)
		for _, row := range rows {
			record, err := getRecord(row, p.opts)

			if err != nil {
				var parseErr *ParseError
				if errors.As(err, &parseErr) {
					p.logger.Printf("Ошибка парсинга счетчика: %v %v", err, row)
				}

				continue
			}

			res = append(res, record)
		}
	}

	return res, nil
}

func getRecord(row []string, opts ParseOptions) (domain.VerificationRecord, error) {
	if row == nil {
		return domain.VerificationRecord{}, errors.New("Row is nil")
	}

	var err error
	// Fields amount without ambient fields
	rowLen := 8
	postfix := ""

	if !rowIsMeter(row) {
		return domain.VerificationRecord{}, errors.New("Not meter")
	}

	if opts.ParseAmbient {
		// fields amount with ambient fields
		rowLen = 11
		postfix = ", температура, влажность, атм. давление"
	}

	if len(row) < rowLen {
		return domain.VerificationRecord{}, ParseError{"Не хватает полей. Ожидаются: номер строки, рег. номер, заводской номер, тип СИ, дата поверки, дата след. поверки, номер установки, поверитель" + postfix}
	}

	res := domain.VerificationRecord{
		RegNum:         strings.TrimSpace(row[1]),
		ManufactureNum: strings.TrimSpace(row[2]),
		MIType:         strings.TrimSpace(row[3]),
		Verifier:       strings.TrimSpace(row[6]),
		TestDeviceNum:  strings.TrimSpace(row[7]),
	}

	res.VerificationDate, err = parseDate(row[4])
	if err != nil {
		return domain.VerificationRecord{}, err
	}

	res.ValidUntil, _ = parseDate(row[5])
	
	if opts.ParseAmbient {
		temp, err := strconv.ParseFloat(row[8], 64)
		if err != nil {
			return domain.VerificationRecord{}, ParseError{"Невозможно преобразовать температуру во float: " + row[8]}
		}

		humidity, err := strconv.ParseFloat(row[9], 64)
		if err != nil {
			return domain.VerificationRecord{}, ParseError{"Невозможно преобразовать влажность во float: " + row[9]}
		}

		pressure, err := strconv.ParseFloat(row[10], 64)
		if err != nil {
			return domain.VerificationRecord{}, ParseError{"Невозможно преобразовать атм.давление во float: " + row[10]}
		}

		res.Temperature = temp
		res.Humidity = humidity
		res.Pressure = pressure
	}

	return res, nil
}

func rowIsMeter(row []string) bool {
	if _, err := strconv.Atoi(row[0]); err != nil {
		return false
	}

	return true
}

func parseDate(dateString string) (time.Time, error) {
	for _, layout := range dataLayouts {
		res, err := time.Parse(layout, dateString)
		if err == nil {
			return res, nil
		}
	}

	return time.Time{}, ParseError{"Некорректная дата: " + dateString}
}
