package service

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"bank/pkg/log"
	"context"
)

// Repository реализует интерфейс репозитория счетов.
type Repository interface {
	CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error)      // Создание счёта
	DeleteAccount(ctx context.Context, accountId uint) error                                  // Удаление счёта
	ListAccounts(ctx context.Context, clientId uint) ([]*models.Account, error)               // Получение списка счетов
	AccountByID(ctx context.Context, accountId uint) (*models.Account, error)                 // Получение счёта по его идентификатору
	AccountByNumber(ctx context.Context, accountNumber string) (*models.Account, error)       // Получение счёта по его номеру
	AccountExistsByNumber(ctx context.Context, accountNumber string) (*models.Account, error) // Получение счёта по его номеру
}

// AuthService реализует интерфейс сервиса пользователей.
type AuthService interface {
	ClientById(ctx context.Context, clientId uint) (client *models.Client, err error) // Получение пользователя по его ID
}

type AccountOption func(*accountService) error

type accountService struct {
	repo Repository
	auth AuthService

	log log.Logger
}

// NewService возвращает имплементацию сервиса для счетов.
func NewService(opts ...AccountOption) (*accountService, error) {
	var svc accountService

	for _, opt := range opts {
		opt(&svc)
	}

	if svc.log == nil {
		return nil, errors.Errorf("no logger provided")
	}

	if svc.repo == nil {
		return nil, errors.Errorf("no repository provided")
	}

	return &svc, nil
}

func WithLogger(logger log.Logger) AccountOption {
	return func(r *accountService) error {
		r.log = logger
		return nil
	}
}

func WithRepository(repo Repository) AccountOption {
	return func(r *accountService) error {
		r.repo = repo
		return nil
	}
}

func WithAuthService(auth AuthService) AccountOption {
	return func(r *accountService) error {
		r.auth = auth
		return nil
	}
}
