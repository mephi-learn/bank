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
	service accountService
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

// WithService добавляет [accountService] в обработчик запросов.
func WithService(svc accountService) HandlerOption {
	return func(o *handler) {
		o.service = svc
	}
}

type accountService interface {
	CreateAccount(ctx context.Context, clientId uint) (account *models.Account, err error) // Создание счёта
	DeleteAccount(ctx context.Context, account *models.Account) error                      // Удаление счёта
	ListAccounts(ctx context.Context, clientId uint) ([]*models.Account, error)            // Список счетов
}

type RouteOption func(*handler)

func (h *handler) WithRouter(mux *http.ServeMux) {
	mux.Handle("POST /api/account/create", auth.Authenticated(h.CreateAccount))
	mux.Handle("PUT /api/account/delete/{account}", auth.Authenticated(h.DeleteAccount))
	mux.Handle("POST /api/account/list", auth.Authenticated(h.ListAccount))

	mux.Handle("POST /api/admin/account/{client}/create", auth.AuthenticatedAdmin(h.CreateAccountAdmin))
	mux.Handle("PUT /api/admin/account/{client}/delete/{account}", auth.AuthenticatedAdmin(h.DeleteAccountAdmin))
	mux.Handle("POST /api/admin/account/{client}/list", auth.AuthenticatedAdmin(h.ListAccountAdmin))
}
