package service

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"context"
)

const createAccountAttempts = 10

// CreateAccount создание счёта.
func (s *accountService) CreateAccount(ctx context.Context, clientId uint) (*models.Account, error) {
	// Проверяем, существует ли пользователь
	if _, err := s.auth.ClientById(ctx, clientId); err != nil {
		s.log.ErrorContext(ctx, "error search user", "error", err, "client id", clientId)
		return nil, errors.Errorf("failed to search user: %w", err)
	}

	// Пробуем создать номер счёта, но чтобы не зависнуть, например, в случае недоступности БД, ограничиваем количество попыток
	attempt := 0
	var accountNumber string
	for accountNumber == "" && attempt < createAccountAttempts {
		newNumber := createAccountNumber()
		if account, err := s.repo.AccountExistsByNumber(ctx, newNumber); account == nil && err == nil {
			accountNumber = newNumber
		}
		attempt++
	}
	if accountNumber == "" {
		s.log.ErrorContext(ctx, "error create account number", "attempts", attempt, "client id", clientId)
		return nil, errors.Errorf("failed to create account number")
	}

	// Формируем данные для создания
	accountRepo := models.Account{
		ClientID: clientId,
		Number:   accountNumber,
		Balance:  0,
	}

	// Создаём счёт в репозитории
	account, err := s.repo.CreateAccount(ctx, &accountRepo)
	if err != nil {
		s.log.ErrorContext(ctx, "error create account", "error", err, "client id", clientId)
		return nil, errors.Errorf("cannot create account: %w", err)
	}

	return account, nil
}
