package services

import (
	"context"
	"errors"
	"movie-festival-be/internal/app/models"
	"movie-festival-be/internal/app/repositories"
	"movie-festival-be/package/helper"
	"movie-festival-be/package/logging"
)

type AuthService interface {
	Login(ctx context.Context, request models.LoginRequest) (string, error)
}

type AuthServiceImpl struct {
	authRepo repositories.AuthRepository
}

func NewAuthService(authRepo repositories.AuthRepository) AuthService {
	return &AuthServiceImpl{authRepo: authRepo}
}

func (s *AuthServiceImpl) Login(ctx context.Context, request models.LoginRequest) (string, error) {
	user, err := s.authRepo.FindUserByEmail(request.Email)
	if err != nil {
		logging.LogError(ctx, "Error fetching user: %v", err)
		return "", errors.New("invalid email or password")
	}

	err = helper.CompareHashAndPassword(user.Password, request.Password)
	if err != nil {
		logging.LogError(ctx, "Password mismatch: %v", err)
		return "", errors.New("invalid username or password")
	}

	token, err := helper.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		logging.LogError(ctx, "Error generating JWT: %v", err)
		return "", errors.New("internal server error")
	}

	if err := s.authRepo.UpdateLoginStatus(user.ID, true); err != nil {
		logging.LogError(ctx, "Error updating login status: %v", err)
		return "", errors.New("internal server error")
	}

	return token, nil
}
