package app

import (
	"arshin/internal/domain"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type RecordsLoader interface {
	Load(string) ([]domain.VerificationRecord, error)
}

type Sink interface {
	Send(string, []domain.Verification) error
}

type ConfigAdapter interface {
	Load(context.Context) (*domain.Config, error)
}

type app struct {
	configAdapter ConfigAdapter
	logger        *log.Logger
}

func New(configAdapter ConfigAdapter, logger *log.Logger) *app {
	return &app{configAdapter: configAdapter, logger: logger}
}

func (a *app) Handle(source string, recordsLoader RecordsLoader, dest string, sink Sink) error {
	var res []domain.Verification

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	config, err := a.configAdapter.Load(ctx)
	var configErr *domain.ErrInvalidConfig
	if errors.As(err, &configErr) {
		return fmt.Errorf("Не удалось загрузить конфиг: %v", err)
	} else if err != nil {
		return fmt.Errorf("Ошибка приложения. Не удалось загрузить конфиг")
	}

	records, err := recordsLoader.Load(source)
	if err != nil {
		return err
	}

	for _, record := range records {
		verification, err := domain.NewVerification(record, config)
		if err != nil {
			a.logger.Printf("Невозможно передать счетчик %s: %v", record.ManufactureNum, err)
			continue
		}

		res = append(res, verification)
	}

	if err := sink.Send(dest, res); err != nil {
		return err
	}

	return nil
}
