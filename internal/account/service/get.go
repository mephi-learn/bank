package service

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"context"
)

// AccountById получение счёта по его идентификатору.
func (s *accountService) AccountByID(ctx context.Context, accountID uint) (*models.Account, error) {
	// Получаем счёт из хранилища
	account, err := s.repo.AccountByID(ctx, accountID)
	if err != nil {
		s.log.ErrorContext(ctx, "error get account", "error", err, "account id", accountID)
		return nil, errors.Errorf("error get account: %w", err)
	}

	return account, nil
}
