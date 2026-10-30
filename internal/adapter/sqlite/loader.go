package sqlite

import (
	"arshin/internal/domain"
	"context"
	"database/sql"
	"errors"
)

type loader struct {
	db *sql.DB
}

func New(db *sql.DB) *loader {
	return &loader{db: db}
}

func (l *loader) Load(ctx context.Context) (*domain.Config, error) {
	signCipher, stringMode, err := l.loadConfig(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("В базе отсутствуют шифр организации и режим обработки показаний окружающей среды")
	} else if err != nil {
		return nil, err
	}

	ambientMode, ok := ambientMode[stringMode]
	if !ok {
		return nil, errors.New("Не удалось сопоставить режим обратоки показаний окр.среды со значением")
	}

	miTypes, err := l.loadMitypes(ctx)
	if err != nil {
		return nil, err
	}
	if len(miTypes) == 0 {
		return nil, errors.New("Запрос на получение поверяемых типов СИ не вернул результата")
	}

	testDevices, err := l.loadTestDevices(ctx)
	if err != nil {
		return nil, err
	}
	if len(testDevices) == 0 {
		return nil, errors.New("Запрос на получение поверочных установок организации не вернул результата")
	}

	return &domain.Config{
		SignCipher:  signCipher,
		AmbientMode: ambientMode,
		Mitypes:     miTypes,
		TestDevices: testDevices,
	}, nil
}

func (l *loader) loadConfig(ctx context.Context) (signCipher, stringMode string, err error) {
	query := `SELECT signcipher, ambientmode FROM config`

	err = l.db.QueryRowContext(ctx, query).Scan(&signCipher, &stringMode)
	if err != nil {
		return "", "", err
	}

	return signCipher, stringMode, nil
}

func (l *loader) loadMitypes(ctx context.Context) (map[string]string, error) {
	query := `SELECT regnum, name FROM mitypes JOIN methods ON (method = id)`
	res := make(map[string]string)

	rows, err := l.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var miitype string
		var method string
		err = rows.Scan(&miitype, &method)

		if err != nil {
			return nil, err
		}

		res[miitype] = method
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (l *loader) loadTestDevices(ctx context.Context) (map[string]domain.TestDevice, error) {
	query := `SELECT testdevices.regnum AS regnum, t1.regnum AS stopwatchregnum, t1.number AS stopwatchnumber, 
							t2.regnum AS thermometerregnum, t2.number AS thermometernumber, 
							t3.regnum AS thermohygrometerregnum, t3.number AS thermohygrometernumber
							FROM testdevices 
								LEFT JOIN testdevicemi AS t1 ON (stopwatch = t1.id) 
								LEFT JOIN testdevicemi AS t2 ON (thermometer = t2.id) 
								LEFT JOIN testdevicemi AS t3 ON (thermohygrometer = t3.id);`
	res := make(map[string]domain.TestDevice)
	rows, err := l.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var regnum string
		var stopwatchregnum string
		var stopwatchnumber string
		var thermometerregnum string
		var thermometernumber string
		var thermohygrometerregnum string
		var thermohygrometernumber string

		err = rows.Scan(query, &regnum, &stopwatchregnum, &stopwatchnumber, &thermometerregnum, &thermometernumber, &thermohygrometerregnum, &thermohygrometernumber)
		if err != nil {
			return nil, err
		}

		res[regnum] = domain.TestDevice{
			Number: regnum,
			Stopwatch: domain.MeasuringInstrument{
				RegNum:         stopwatchregnum,
				ManufactureNum: stopwatchnumber,
			},
			Thermometer: domain.MeasuringInstrument{
				RegNum:         thermometerregnum,
				ManufactureNum: thermometernumber,
			},
			Thermohygrometer: domain.MeasuringInstrument{
				RegNum:         thermohygrometerregnum,
				ManufactureNum: thermohygrometernumber,
			},
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
