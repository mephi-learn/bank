package service

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"context"
)

// ListAccounts список счетов.
func (s *accountService) ListAccounts(ctx context.Context, clientId uint) ([]*models.Account, error) {
	// Получаем список счётов из хранилища
	accounts, err := s.repo.ListAccounts(ctx, clientId)
	if err != nil {
		s.log.ErrorContext(ctx, "error list accounts", "error", err, "client id", clientId)
		return nil, errors.Errorf("error list accounts: %w", err)
	}

	return accounts, nil
}
