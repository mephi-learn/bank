package controller

import (
	"bank/internal/auth"
	"bank/internal/models"
	"bank/pkg/errors"
	"bank/pkg/log"
	"context"
	"net/http"
)

type handler struct {
	service transactionService
	log     log.Logger
}

type HandlerOption func(*handler)

func NewHandler(opts ...HandlerOption) (*handler, error) {
	h := &handler{}

	for _, opt := range opts {
		opt(h)
	}

	if h.log == nil {
		return nil, errors.New("logger is missing")
	}

	return h, nil
}

func WithLogger(logger log.Logger) HandlerOption {
	return func(o *handler) {
		o.log = logger
	}
}

// WithService добавляет [transactionService] в обработчик запросов.
func WithService(svc transactionService) HandlerOption {
	return func(o *handler) {
		o.service = svc
	}
}

type transactionService interface {
	CardByCard(ctx context.Context, fromCard string, toCard string, amount float64) ([]*models.Transaction, error) // Перевод между картами
	Deposit(ctx context.Context, cardNum string, amount float64) (*models.Transaction, error)                      // Пополнение баланса
}

type RouteOption func(*handler)

func (h *handler) WithRouter(mux *http.ServeMux) {
	mux.Handle("POST /api/transaction/{from_card_num}/{to_card_num}/{amount}", auth.Authenticated(h.CardToCard)) // Перевод между картами
	mux.Handle("POST /api/transaction/deposit/{card_num}/{amount}", auth.Authenticated(h.Deposit))               // Пополнение карта
}
