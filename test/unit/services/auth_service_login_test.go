package services_test

import (
	"movie-festival-be/internal/app/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginService(t *testing.T) {
	request := models.LoginRequest{
		Email:    "admin@movie-festival.com",
		Password: "Symantec2121",
	}

	token, err := authService.Login(ctx, request)
	assert.Nil(t, err)
	assert.NotEqual(t, "", token)
}
