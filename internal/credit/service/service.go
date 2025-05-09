package service

import (
	"bank/pkg/errors"
	"bank/pkg/log"
)

type CreditOption func(*creditService) error

type creditService struct {
	log log.Logger
}

// NewService возвращает имплементацию сервиса для счетов.
func NewService(opts ...CreditOption) (*creditService, error) {
	var svc creditService

	for _, opt := range opts {
		opt(&svc)
	}

	if svc.log == nil {
		return nil, errors.Errorf("no logger provided")
	}

	return &svc, nil
}

func WithLogger(logger log.Logger) CreditOption {
	return func(r *creditService) error {
		r.log = logger
		return nil
	}
}
