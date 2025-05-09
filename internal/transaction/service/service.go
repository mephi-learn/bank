package service

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"bank/pkg/log"
	"context"
)

// Repository реализует интерфейс репозитория счетов.
type Repository interface {
	TransferBetweenAccounts(ctx context.Context, accountFrom uint, accountTo uint, amount float64) ([]*models.Transaction, error)
	Deposit(ctx context.Context, account uint, amount float64) (*models.Transaction, error)
}

// CardService реализует интерфейс сервиса карт.
type CardService interface {
	ValidateCardByNum(ctx context.Context, cardNum string) (*models.Card, error) // Проверка карты по её номеру
}

// AccountService реализует интерфейс сервиса счетов.
type AccountService interface {
	AccountByID(ctx context.Context, accountID uint) (*models.Account, error) // Получение счёта по его ID
}

type TransactionOption func(*transactionService) error

type transactionService struct {
	repo    Repository
	account AccountService
	card    CardService

	log log.Logger
}

// NewService возвращает имплементацию сервиса для счетов.
func NewService(opts ...TransactionOption) (*transactionService, error) {
	var svc transactionService

	for _, opt := range opts {
		opt(&svc)
	}

	if svc.log == nil {
		return nil, errors.Errorf("no logger provided")
	}

	if svc.repo == nil {
		return nil, errors.Errorf("no repository provided")
	}

	if svc.account == nil {
		return nil, errors.Errorf("no account provided")
	}

	return &svc, nil
}

func WithLogger(logger log.Logger) TransactionOption {
	return func(r *transactionService) error {
		r.log = logger
		return nil
	}
}

func WithRepository(repo Repository) TransactionOption {
	return func(r *transactionService) error {
		r.repo = repo
		return nil
	}
}

func WithAccountService(account AccountService) TransactionOption {
	return func(r *transactionService) error {
		r.account = account
		return nil
	}
}

func WithCardService(card CardService) TransactionOption {
	return func(r *transactionService) error {
		r.card = card
		return nil
	}
}
