package service

import (
	"bank/internal/models"
	"bank/pkg/errors"
	"context"
	"time"
)

func (h *creditService) CreateCredit(ctx context.Context) (account *models.Credit, err error) {
	credit, err := h.CreateCredit(ctx)
	if err != nil {
		return nil, errors.Errorf("failed create credit: %w", err)
	}

	return credit, nil
}

func (h *creditService) PayCredit(ctx context.Context, payDate time.Time) (*models.Payment, error) {
	pay, err := h.PayCredit(ctx, payDate)
	if err != nil {
		return nil, errors.Errorf("failed pay credit: %w", err)
	}

	return pay, nil
}

func (h *creditService) ListCards(ctx context.Context) ([]*models.Credit, error) {
	creditList, err := h.ListCards(ctx)
	if err != nil {
		return nil, errors.Errorf("failed list credit: %w", err)
	}

	return creditList, nil
}
