package service

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"context"
)

// CardById получение карты по его идентификатору.
func (s *cardService) CardById(ctx context.Context, cardID uint) (*models.Card, error) {
	// Получаем карту из хранилища
	card, err := s.repo.CardById(ctx, cardID)
	if err != nil {
		s.log.ErrorContext(ctx, "error get card", "error", err, "card id", cardID)
		return nil, errors.Errorf("error get card: %w", err)
	}

	return card, nil
}
