package service

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"context"
)

// DeleteCard удаление карты.
func (s *cardService) DeleteCard(ctx context.Context, card *models.Card) error {
	// Проверяем, существует ли карта
	cardRepo, err := s.repo.CardById(ctx, card.ID)
	if err != nil {
		s.log.ErrorContext(ctx, "error check card", "error", err, "card id", card.ID)
		return errors.Errorf("failed to search card: %w", err)
	}

	// Проверяем, что клиент указанный в запросе на удаление и клиент в карте, совпадают
	if card.ClientID != cardRepo.ClientID {
		s.log.ErrorContext(ctx, "delete card another client", "error", err, "card id", card.ID)
		return errors.Errorf("delete card another client: %w", err)
	}

	// Удаляем счёт из хранилища
	if err = s.repo.DeleteCard(ctx, card.ID); err != nil {
		s.log.ErrorContext(ctx, "error delete card", "error", err, "card id", card.ID)
		return errors.Errorf("cannot delete card: %w", err)
	}

	return nil
}

// DeleteCardByAdmin удаление любой карты.
func (s *cardService) DeleteCardByAdmin(ctx context.Context, card *models.Card) error {
	// Проверяем, существует ли карта
	if _, err := s.repo.CardById(ctx, card.ID); err != nil {
		s.log.ErrorContext(ctx, "error check card", "error", err, "card id", card.ID)
		return errors.Errorf("failed to search card: %w", err)
	}

	// Удаляем счёт из хранилища
	if err := s.repo.DeleteCard(ctx, card.ID); err != nil {
		s.log.ErrorContext(ctx, "error delete card", "error", err, "card id", card.ID)
		return errors.Errorf("cannot delete card: %w", err)
	}

	return nil
}
