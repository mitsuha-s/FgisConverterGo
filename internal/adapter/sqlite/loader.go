package sqlite

import (
	"arshin/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type repo struct {
	db *sql.DB
}

func New(db *sql.DB) *repo {
	return &repo{db: db}
}

func (r *repo) GetConfig(ctx context.Context) (*domain.Config, error) {
	signCipher, stringMode, err := r.loadConfig(ctx)
	if err != nil {
		return nil, err
	}

	ambientMode, ok := ambientMode[stringMode]
	if !ok {
		return nil, errors.New("Не удалось сопоставить режим обработки показаний окр.среды со значением")
	}

	miTypes, err := r.loadMitypes(ctx)
	if err != nil {
		return nil, err
	}

	methods, err := r.loadMethods(ctx)
	if err != nil {
		return nil, err
	}

	testDevices, err := r.loadTestDevices(ctx)
	if err != nil {
		return nil, err
	}

	return domain.NewConfig(signCipher, miTypes, methods, testDevices, ambientMode)
}

func (r *repo) loadConfig(ctx context.Context) (signCipher, stringMode string, err error) {
	query := `SELECT signcipher, ambientmode FROM config`

	var sign sql.NullString
	var mode sql.NullString

	err = r.db.QueryRowContext(ctx, query).Scan(&sign, &mode)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", errors.New("В базе отсутствуют шифр организации и режим обработки показаний окружающей среды")
	} else if err != nil {
		return "", "", fmt.Errorf("Не удалось загрузить загрузить шифр организации и режим обработки показаний окружающей среды: %w", err)
	}

	return sign.String, mode.String, nil
}

func (r *repo) loadMitypes(ctx context.Context) (map[string]string, error) {
	query := `SELECT regnum, name FROM mitypes JOIN methods ON (method = id)`
	res := make(map[string]string)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Не удалось загрузить поверяемые типы СИ: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var miitype string
		var method string
		err = rows.Scan(&miitype, &method)

		if err != nil {
			return nil, fmt.Errorf("Не удалось загрузить поверяемые типы СИ: %w", err)
		}

		res[miitype] = method
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Не удалось загрузить поверяемые типы СИ: %w", err)
	}

	return res, nil
}

func (r *repo) loadTestDevices(ctx context.Context) ([]domain.TestDevice, error) {
	query := `SELECT testdevicemi.regnum, testdevicemi.number, testdevices.regnum 
				FROM testdevicemi RIGHT JOIN testdevices ON (testdevice = id)`

	res := make([]domain.TestDevice, 0)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Не удалось загрузить поверочные установки: %w", err)
	}

	defer rows.Close()

	instrumentsMap := make(map[string][]domain.MeasuringInstrument)

	for rows.Next() {
		var regnum string
		var instrumentregnum string
		var instrumentnumber string

		err = rows.Scan(&instrumentregnum, &instrumentnumber, &regnum)
		if err != nil {
			return nil, fmt.Errorf("Не удалось загрузить поверочные установки: %w", err)
		}
		instrumentsMap[regnum] = nil

		instrument, err := domain.NewMeasuringInstrument(instrumentregnum, instrumentnumber)
		if err != nil {
			return nil, fmt.Errorf("Не удалось загрузить доп.средство измерения поверочной установки %s %s: %w", instrumentregnum, instrumentnumber, err)
		}

		instrumentsMap[regnum] = append(instrumentsMap[regnum], instrument)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Не удалось загрузить поверочные установки: %w", err)
	}

	for testdeviceNum, instruments := range instrumentsMap {
		testdevice, err := domain.NewTestDevice(testdeviceNum, instruments)
		if err != nil {
			return nil, fmt.Errorf("Не удалось загрузить поверочную установку %s: %w", testdeviceNum, err)
		}

		res = append(res, testdevice)
	}

	return res, nil
}

func (r *repo) loadMethods(ctx context.Context) ([]string, error) {
	var res []string

	query := `SELECT name FROM methods`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Не удалось загрузить список разрешенных методик: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var method string

		err := rows.Scan(&method)
		if err != nil {
			return nil, fmt.Errorf("Не удалось загрузить список разрешенных методик: %w", err)
		}

		res = append(res, method)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Не удалось загрузить список разрешенных методик: %w", err)
	}

	return res, nil
}
