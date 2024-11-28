package handlers

import (
	"encoding/json"
	"fmt"
	"movie-festival-be/internal/app/models"
	"movie-festival-be/internal/app/services"
	"movie-festival-be/package/helper"
	"movie-festival-be/package/logging"
	"movie-festival-be/package/middleware"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Authhandler struct {
	authService services.AuthService
}

func NewAuthhandler(authService services.AuthService) *Authhandler {
	return &Authhandler{authService: authService}
}

func (h *Authhandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Println(r.Body)
	var loginRequest models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		logging.LogError(ctx, "Error decoding request body: %v", err)
		middleware.WriteResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(loginRequest); err != nil {
		middleware.WriteResponse(w, http.StatusBadRequest, helper.GetMessageValidator(validate, err), nil)
		return
	}

	token, err := h.authService.Login(ctx, loginRequest)
	if err != nil {
		logging.LogError(ctx, "Login error: %v", err)
		middleware.WriteResponse(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	middleware.WriteResponse(w, http.StatusOK, "", map[string]string{"token": token})
}
