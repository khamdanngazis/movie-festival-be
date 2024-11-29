package handlers

import (
	"movie-festival-be/internal/app/entities"
	"movie-festival-be/internal/app/services"
	"movie-festival-be/package/middleware"
	"net/http"
)

type StatsHandler struct {
	statsService services.StatsService
	authService  services.AuthService
}

func NewStatsHandler(statsService services.StatsService, authService services.AuthService) *StatsHandler {
	return &StatsHandler{statsService: statsService, authService: authService}
}

func (h *StatsHandler) GetAdminStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok || user.Role != entities.RoleAdmin {
		middleware.WriteResponse(w, http.StatusForbidden, "Access denied", nil)
		return
	}

	if !h.authService.IsLogin(ctx, user.UserID) {
		middleware.WriteResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	stats, err := h.statsService.GetAdminStats(ctx)
	if err != nil {
		middleware.WriteResponse(w, http.StatusInternalServerError, "Failed to fetch statistics", nil)
		return
	}

	middleware.WriteResponse(w, http.StatusOK, "Statistics fetched successfully", stats)
}
