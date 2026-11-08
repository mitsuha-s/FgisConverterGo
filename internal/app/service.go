package app

import (
	"arshin/internal/domain"
	"context"
	"errors"
	"fmt"
	"log"
	// "time"
)

type RecordsLoader interface {
	Load(string) ([]domain.VerificationRecord, error)
}

type Deliverer interface {
	Deliver(string, []domain.Verification) error
}

type DBrepo interface {
	GetConfig(context.Context) (*domain.Config, error)
}

type Service struct {
	configAdapter DBrepo
	logger        *log.Logger
}

func New(configAdapter DBrepo, logger *log.Logger) *Service {
	return &Service{configAdapter: configAdapter, logger: logger}
}

func (a *Service) Handle(source string, recordsLoader RecordsLoader, dest string, sink Deliverer) error {
	var res []domain.Verification
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	// defer cancel()

	config, err := a.configAdapter.GetConfig(context.Background())
	var configErr domain.ErrInvalidData
	if errors.As(err, &configErr) {
		return fmt.Errorf("Не удалось загрузить конфиг: %w", err)
	} else if err != nil {
		return err
	}

	records, err := recordsLoader.Load(source)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return fmt.Errorf("%w. Проверьте порядок и количество столбов или формат файла", domain.ErrNoRecords)
	}

	for _, record := range records {
		verification, err := domain.NewVerification(record, config)
		if err != nil {
			a.logger.Printf("Невозможно передать счетчик %s: %v", record.ManufactureNum, err)
			continue
		}

		res = append(res, verification)
	}

	if len(res) != 0 {
		if err := sink.Deliver(dest, res); err != nil {
			return err
		}
	}

	return nil
}
