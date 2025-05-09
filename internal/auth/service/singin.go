package service

import (
	"bank/internal/auth"
	"bank/internal/models"
	"bank/pkg/errors"
	"context"
	"encoding/json"
)

func (h *authService) SignIn(ctx context.Context, userData *models.Client) (string, error) {
	userData.Password = generatePasswordHash(userData.Username, userData.Password)
	user, err := h.repo.GetByUsernameAndPassword(ctx, userData)
	if err != nil {
		return "", err
	}

	marshaledUser, err := json.Marshal(user)
	if err != nil {
		return "", errors.Errorf("failed marshal user data: %w", err)
	}

	signedToken, err := auth.GenerateJWTToken(string(marshaledUser))
	if err != nil {
		return "", errors.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
