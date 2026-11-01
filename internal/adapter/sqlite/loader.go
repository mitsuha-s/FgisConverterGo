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

	testDevices, err := l.loadTestDevices(ctx)
	if err != nil {
		return nil, err
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

	var sign sql.NullString
	var mode sql.NullString

	err = l.db.QueryRowContext(ctx, query).Scan(&sign, &mode)
	if err != nil {
		return "", "", err
	}

	return sign.String, mode.String, nil
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
		var regnum sql.NullString
		var stopwatchregnum sql.NullString
		var stopwatchnumber sql.NullString
		var thermometerregnum sql.NullString
		var thermometernumber sql.NullString
		var thermohygrometerregnum sql.NullString
		var thermohygrometernumber sql.NullString

		err = rows.Scan(&regnum, &stopwatchregnum, &stopwatchnumber, &thermometerregnum, &thermometernumber, &thermohygrometerregnum, &thermohygrometernumber)
		if err != nil {
			return nil, err
		}

		res[regnum.String] = domain.TestDevice{
			Number: regnum.String,
			Stopwatch: domain.MeasuringInstrument{
				RegNum:         stopwatchregnum.String,
				ManufactureNum: stopwatchnumber.String,
			},
			Thermometer: domain.MeasuringInstrument{
				RegNum:         thermometerregnum.String,
				ManufactureNum: thermometernumber.String,
			},
			Thermohygrometer: domain.MeasuringInstrument{
				RegNum:         thermohygrometerregnum.String,
				ManufactureNum: thermohygrometernumber.String,
			},
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
