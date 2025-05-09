package service

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"context"
)

// ListCards список карт.
func (s *cardService) ListCards(ctx context.Context, clientID uint) ([]*models.Card, error) {
	// Получаем список карт из хранилища
	accounts, err := s.repo.ListCards(ctx, clientID)
	if err != nil {
		s.log.ErrorContext(ctx, "error list cards", "error", err, "card id", clientID)
		return nil, errors.Errorf("error list cards: %w", err)
	}

	return accounts, nil
}
