package service

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"context"
)

// DeleteAccount удаление счёта.
func (s *accountService) DeleteAccount(ctx context.Context, account *models.Account) error {
	// Проверяем, существует ли счёт
	accountRepo, err := s.repo.AccountByID(ctx, account.ID)
	if err != nil {
		s.log.ErrorContext(ctx, "error check account", "error", err, "account id", account.ID)
		return errors.Errorf("failed to search account: %w", err)
	}

	// Проверяем, что клиент указанный в запросе на удаление и клиент в счёте, совпадают
	if account.ClientID != accountRepo.ClientID {
		s.log.ErrorContext(ctx, "delete account another client", "error", err, "account id", account.ID)
		return errors.Errorf("delete account another client: %w", err)
	}

	// Проверяем, что на счету нет денег
	if accountRepo.Balance > 0 {
		s.log.ErrorContext(ctx, "account has money", "account_id", account.ID)
		return errors.Errorf("account has money: %d", accountRepo.ID)
	}

	// Удаляем счёт из хранилища
	if err = s.repo.DeleteAccount(ctx, account.ID); err != nil {
		s.log.ErrorContext(ctx, "error delete account", "error", err, "account id", account.ID)
		return errors.Errorf("cannot delete account: %w", err)
	}

	return nil
}
