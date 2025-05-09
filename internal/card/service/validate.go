package service

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"context"
)

// ValidateCardByNum валидация карты по её номеру.
func (s *cardService) ValidateCardByNum(ctx context.Context, cardNum string) (*models.Card, error) {
	if !cardNumberIsValid(cardNum) {
		return nil, errors.New("invalid card number: %s" + cardNum)
	}

	// Получаем карту из хранилища
	card, err := s.repo.CardByNumber(ctx, cardNum)
	if err != nil {
		s.log.ErrorContext(ctx, "error get card by number", "error", err, "card number", cardNum)
		return nil, errors.Errorf("error get card by number: %w", err)
	}

	if card.Expired() {
		return nil, errors.New("card expire: %s" + cardNum)
	}

	return card, nil
}
