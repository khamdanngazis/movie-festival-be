package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogoutService(t *testing.T) {

	user, _ := authRepo.FindUserByEmail("admin@movie-festival.com")

	err := authService.Logout(ctx, user.ID)
	assert.Nil(t, err)
	user, err = authRepo.FindUserByEmail("admin@movie-festival.com")
	assert.Nil(t, err)
	assert.True(t, !user.LoggedIn)
}
