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
	CreateCard(ctx context.Context, accountId uint) (account *models.Card, err error) // Создание карты
	CreateCardByAdmin(ctx context.Context, accountID uint) (*models.Card, error)      // Создание карты администратором
	DeleteCard(ctx context.Context, account *models.Card) error                       // Удаление карты
	DeleteCardByAdmin(ctx context.Context, account *models.Card) error                // Удаление карты администратором
	ListCards(ctx context.Context, accountId uint) ([]*models.Card, error)            // Список карт
}

type RouteOption func(*handler)

func (h *handler) WithRouter(mux *http.ServeMux) {
	mux.Handle("POST /api/card/{account_id}/create", auth.Authenticated(h.CreateCard)) // Создание карты пользователем (себе)
	mux.Handle("PUT /api/card/delete/{card_id}", auth.Authenticated(h.DeleteCard))     // Удаление карты пользователем (своей)
	mux.Handle("POST /api/card/list", auth.Authenticated(h.ListCards))                 // Список карт пользователя (своих)

	mux.Handle("POST /api/admin/card/{account_id}/create", auth.AuthenticatedAdmin(h.CreateCardByAdmin)) // Создание карты администратором (чужой)
	mux.Handle("PUT /api/admin/card/delete/{card_id}", auth.AuthenticatedAdmin(h.DeleteCardByAdmin))     // Удаление карты администратором (чужой)
	mux.Handle("POST /api/admin/card/{client_id}/list", auth.AuthenticatedAdmin(h.ListCardsByAdmin))     // Список карт пользователя под администратором (чужих)
}
