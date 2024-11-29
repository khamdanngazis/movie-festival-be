package services

import (
	"context"
	"errors"
	"movie-festival-be/internal/app/entities"
	"movie-festival-be/internal/app/models"
	"movie-festival-be/internal/app/repositories"
	"movie-festival-be/package/helper"
	"movie-festival-be/package/logging"
)

type AuthService interface {
	Login(ctx context.Context, request models.LoginRequest) (string, error)
	IsLogin(ctx context.Context, userID uint) bool
	Register(ctx context.Context, request models.RegisterRequest) error
	Logout(ctx context.Context, userID uint) error
	ExtendToken(ctx context.Context, userID uint) (string, error)
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

func (s *AuthServiceImpl) IsLogin(ctx context.Context, userID uint) bool {
	user, err := s.authRepo.FindUserByID(userID)
	if err != nil {
		logging.LogError(ctx, "Failed to find login status: %v", err)
		return false
	}
	return user.LoggedIn
}

func (s *AuthServiceImpl) Register(ctx context.Context, request models.RegisterRequest) error {
	existingUser, _ := s.authRepo.FindUserByEmail(request.Email)
	if existingUser != nil {
		return errors.New("email is already registered")
	}

	hashedPassword, err := helper.HashPassword(request.Password)
	if err != nil {
		logging.LogError(ctx, "Error hashing password: %v", err)
		return errors.New("internal server error")
	}

	user := &entities.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
		Role:     entities.RoleUser,
	}

	if err := s.authRepo.SaveUser(user); err != nil {
		logging.LogError(ctx, "Error saving user: %v", err)
		return errors.New("internal server error")
	}

	return nil
}

func (s *AuthServiceImpl) Logout(ctx context.Context, userID uint) error {
	if err := s.authRepo.UpdateLoginStatus(userID, false); err != nil {
		logging.LogError(ctx, "Failed to update login status: %v", err)
		return errors.New("failed to log out")
	}
	return nil
}

func (s *AuthServiceImpl) ExtendToken(ctx context.Context, id uint) (string, error) {
	user, err := s.authRepo.FindUserByID(id)
	if err != nil {
		logging.LogError(ctx, "Error fetching user: %v", err)
		return "", errors.New("invalid email or password")
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
