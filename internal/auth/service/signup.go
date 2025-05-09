package service

import (
	"bank/internal/models"
	"context"
)

func (h *authService) SignUp(ctx context.Context, userData *models.Client) (uint, error) {
	userData.Password = generatePasswordHash(userData.Username, userData.Password)
	return h.repo.Create(ctx, userData)
}
