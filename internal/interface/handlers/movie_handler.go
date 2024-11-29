package handlers

import (
	"encoding/json"
	"movie-festival-be/internal/app/entities"
	"movie-festival-be/internal/app/models"
	"movie-festival-be/internal/app/services"
	"movie-festival-be/package/helper"
	"movie-festival-be/package/logging"
	"movie-festival-be/package/middleware"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type MoviesHandler struct {
	movieService services.MovieService
	authService  services.AuthService
}

func NewMoviesHandler(movieService services.MovieService, authService services.AuthService) *MoviesHandler {
	return &MoviesHandler{movieService: movieService, authService: authService}
}

func (h *MoviesHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
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
	var request models.SaveMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logging.LogError(ctx, "Error decoding request body: %v", err)
		middleware.WriteResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		middleware.WriteResponse(w, http.StatusBadRequest, helper.GetMessageValidator(validate, err), nil)
		return
	}

	movie, err := h.movieService.CreateMovie(ctx, request)

	if err != nil {
		logging.LogError(ctx, "Error creating movie: %v", err)
		middleware.WriteResponse(w, http.StatusInternalServerError, "Failed to create movie", nil)
		return
	}

	middleware.WriteResponse(w, http.StatusCreated, "Movie created successfully", movie)
}

func (h *MoviesHandler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
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

	vars := mux.Vars(r)
	movieID, err := strconv.Atoi(vars["id"])
	if err != nil || movieID <= 0 {
		middleware.WriteResponse(w, http.StatusBadRequest, "Invalid movie ID", nil)
		return
	}

	var request models.SaveMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logging.LogError(ctx, "Error decoding request body: %v", err)
		middleware.WriteResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		middleware.WriteResponse(w, http.StatusBadRequest, helper.GetMessageValidator(validate, err), nil)
		return
	}

	movie, err := h.movieService.UpdateMovie(ctx, uint(movieID), request)
	if err != nil {
		logging.LogError(ctx, "Error updating movie: %v", err)
		middleware.WriteResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	middleware.WriteResponse(w, http.StatusOK, "Movie updated successfully", movie)
}

func (h *MoviesHandler) ListMovies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	movies, total, err := h.movieService.ListMovies(ctx, page, limit)
	if err != nil {
		middleware.WriteResponse(w, http.StatusInternalServerError, "Failed to fetch movies", nil)
		return
	}

	middleware.WriteResponse(w, http.StatusOK, "Movies fetched successfully", map[string]interface{}{
		"movies": movies,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

func (h *MoviesHandler) SearchMovies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query().Get("query")
	if query == "" {
		middleware.WriteResponse(w, http.StatusBadRequest, "Query parameter is required", nil)
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1 // Default to page 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // Default to 10 results per page
	}

	movies, total, err := h.movieService.SearchMovies(ctx, query, page, limit)
	if err != nil {
		logging.LogError(ctx, "Error searching movies: %v", err)
		middleware.WriteResponse(w, http.StatusInternalServerError, "Failed to search movies", nil)
		return
	}

	response := map[string]interface{}{
		"movies": movies,
		"total":  total,
		"page":   page,
		"limit":  limit,
	}

	middleware.WriteResponse(w, http.StatusOK, "Movies fetched successfully", response)
}

func (h *MoviesHandler) TrackView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	movieID, err := strconv.Atoi(vars["id"])
	if err != nil || movieID <= 0 {
		middleware.WriteResponse(w, http.StatusBadRequest, "Invalid movie ID", nil)
		return
	}

	ipAddress := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ipAddress = forwarded
	}

	var userID *uint
	user, _ := middleware.GetUserFromContext(r.Context())
	if user != nil {
		userID = &user.UserID
	}

	err = h.movieService.TrackView(ctx, uint(movieID), ipAddress, userID)
	if err != nil {
		middleware.WriteResponse(w, http.StatusInternalServerError, "Failed to track viewership", nil)
		return
	}

	middleware.WriteResponse(w, http.StatusOK, "View tracked successfully", nil)
}
