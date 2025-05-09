package controller

import (
	"bank/internal/auth"
	"bank/internal/models"
	"bank/pkg/errors"
	"bank/pkg/log"
	"context"
	"net/http"
	"time"
)

type handler struct {
	service cardService
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

// WithService добавляет [cardService] в обработчик запросов.
func WithService(svc cardService) HandlerOption {
	return func(o *handler) {
		o.service = svc
	}
}

type cardService interface {
	CreateCredit(ctx context.Context) (account *models.Credit, err error)      // Создание кредита и расчёт графика платежей
	PayCredit(ctx context.Context, payDate time.Time) (*models.Payment, error) // Погашение кредита на указанную дату
	ListCards(ctx context.Context) ([]*models.Credit, error)                   // Список кредитов
}

type RouteOption func(*handler)

func (h *handler) WithRouter(mux *http.ServeMux) {
	mux.Handle("POST /api/creadit/create", auth.Authenticated(h.CreateCredit)) // Создание кредита
	mux.Handle("PUT /api/credit/pay/{date}", auth.Authenticated(h.PayCredit))  // Оплата кредита
	mux.Handle("POST /api/credit/list", auth.Authenticated(h.ListCredit))      // Список кредитов
}
