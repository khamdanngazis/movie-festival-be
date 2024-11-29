package handlers

import (
	"movie-festival-be/internal/app/entities"
	"movie-festival-be/internal/app/services"
	"movie-festival-be/package/middleware"
	"net/http"
)

type ReportHandler struct {
	reportService services.ReportService
	authService   services.AuthService
}

func NewReportHandler(reportService services.ReportService, authService services.AuthService) *ReportHandler {
	return &ReportHandler{reportService: reportService, authService: authService}
}

func (h *ReportHandler) GetReportViews(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user.Role != entities.RoleAdmin {
		middleware.WriteResponse(w, http.StatusForbidden, "Access denied", nil)
		return
	}

	if !h.authService.IsLogin(ctx, user.UserID) {
		middleware.WriteResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	report, err := h.reportService.GetReportViews(ctx)
	if err != nil {
		middleware.WriteResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	middleware.WriteResponse(w, http.StatusOK, "Report fetched successfully", report)
}
