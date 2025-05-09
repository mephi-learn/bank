package service

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"bank/pkg/log"
	"context"
)

// Repository реализует интерфейс репозитория счетов.
type Repository interface {
	CreateCard(ctx context.Context, card *models.Card) (*models.Card, error)         // Создание карты
	DeleteCard(ctx context.Context, cardID uint) error                               // Удаление карты
	ListCards(ctx context.Context, clientID uint) ([]*models.Card, error)            // Получение списка карт
	CardById(ctx context.Context, cardID uint) (*models.Card, error)                 // Получение карты по её идентификатору
	CardByNumber(ctx context.Context, cardNumber string) (*models.Card, error)       // Получение карты по её номеру
	CardExistsByNumber(ctx context.Context, cardNumber string) (*models.Card, error) // Проверка наличия карты по её номеру
}

// AuthService реализует интерфейс сервиса пользователей.
type AuthService interface {
	ClientById(ctx context.Context, clientId uint) (client *models.Client, err error) // Получение пользователя по его ID
}

// AccountService реализует интерфейс сервиса счетов.
type AccountService interface {
	AccountByID(ctx context.Context, accountID uint) (*models.Account, error) // Получение счёта по его ID
}

type AccountOption func(*cardService) error

type cardService struct {
	repo    Repository
	auth    AuthService
	account AccountService

	log log.Logger
}

// NewService возвращает имплементацию сервиса для счетов.
func NewService(opts ...AccountOption) (*cardService, error) {
	var svc cardService

	for _, opt := range opts {
		opt(&svc)
	}

	if svc.log == nil {
		return nil, errors.Errorf("no logger provided")
	}

	if svc.repo == nil {
		return nil, errors.Errorf("no repository provided")
	}

	if svc.auth == nil {
		return nil, errors.Errorf("no auth provided")
	}

	if svc.account == nil {
		return nil, errors.Errorf("no account provided")
	}

	return &svc, nil
}

func WithLogger(logger log.Logger) AccountOption {
	return func(r *cardService) error {
		r.log = logger
		return nil
	}
}

func WithRepository(repo Repository) AccountOption {
	return func(r *cardService) error {
		r.repo = repo
		return nil
	}
}

func WithAuthService(auth AuthService) AccountOption {
	return func(r *cardService) error {
		r.auth = auth
		return nil
	}
}

func WithAccountService(account AccountService) AccountOption {
	return func(r *cardService) error {
		r.account = account
		return nil
	}
}
