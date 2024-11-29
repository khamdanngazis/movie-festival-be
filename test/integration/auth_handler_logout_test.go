package handlers_test

import (
	"movie-festival-be/internal/app/models"
	"movie-festival-be/package/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogoutHandler(t *testing.T) {
	request := models.LoginRequest{
		Email:    "admin@movie-festival.com",
		Password: "Symantec2121",
	}

	token, _ := authService.Login(ctx, request)

	httpRouter.POST("/auth/logout", authHandler.Logout)

	req, err := http.NewRequest("POST", "/auth/logout", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	httpRouter.Mux().Use(middleware.AuthMiddleware)
	httpRouter.Mux().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

}
