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

func (h *Authhandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var registerRequest models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		logging.LogError(ctx, "Error decoding request body: %v", err)
		middleware.WriteResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(registerRequest); err != nil {
		middleware.WriteResponse(w, http.StatusBadRequest, helper.GetMessageValidator(validate, err), nil)
		return
	}

	if err := h.authService.Register(ctx, registerRequest); err != nil {
		logging.LogError(ctx, "Register error: %v", err)
		middleware.WriteResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	middleware.WriteResponse(w, http.StatusCreated, "User registered successfully", nil)
}

func (h *Authhandler) Logout(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.WriteResponse(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userID := userClaims.UserID
	ctx := r.Context()
	if err := h.authService.Logout(ctx, userID); err != nil {
		middleware.WriteResponse(w, http.StatusInternalServerError, "Failed to log out", nil)
		return
	}

	middleware.WriteResponse(w, http.StatusOK, "Logout successful", nil)
}

func (h *Authhandler) ExtendToken(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		middleware.WriteResponse(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	ctx := r.Context()

	if !h.authService.IsLogin(ctx, user.UserID) {
		middleware.WriteResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	token, err := h.authService.ExtendToken(ctx, user.UserID)
	if err != nil {
		logging.LogError(ctx, "Login error: %v", err)
		middleware.WriteResponse(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	middleware.WriteResponse(w, http.StatusOK, "", map[string]string{"token": token})
}
