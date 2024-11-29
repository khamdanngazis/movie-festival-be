package services_test

import (
	"movie-festival-be/internal/app/entities"
	"movie-festival-be/internal/app/models"
	"movie-festival-be/package/helper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterService(t *testing.T) {
	request := models.RegisterRequest{
		Name:     "Joe",
		Email:    "joe@movie-festival.com",
		Password: "Symantec2121",
	}

	err := authService.Register(ctx, request)
	assert.Nil(t, err)
	user, err := authRepo.FindUserByEmail(request.Email)
	assert.Nil(t, err)
	assert.Equal(t, user.Name, request.Name)
	assert.Equal(t, user.Email, request.Email)
	assert.Equal(t, user.Role, entities.RoleUser)
	err = helper.CompareHashAndPassword(user.Password, request.Password)
	assert.Nil(t, err)
}
