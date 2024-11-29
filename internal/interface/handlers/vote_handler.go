package handlers

import (
	"movie-festival-be/internal/app/entities"
	"movie-festival-be/internal/app/services"
	"movie-festival-be/package/middleware"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type VoteHandler struct {
	voteServie  services.VoteService
	authService services.AuthService
}

func NewVoteHandler(voteServie services.VoteService, authService services.AuthService) *VoteHandler {
	return &VoteHandler{voteServie: voteServie, authService: authService}
}

func (h *VoteHandler) VoteMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		middleware.WriteResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
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

	err = h.voteServie.VoteMovie(ctx, uint(movieID), user.UserID)
	if err != nil {
		if err.Error() == "user has already voted for this movie" {
			middleware.WriteResponse(w, http.StatusConflict, err.Error(), nil)
		} else {
			middleware.WriteResponse(w, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	middleware.WriteResponse(w, http.StatusOK, "Vote registered successfully", nil)
}

func (h *VoteHandler) Unvote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		middleware.WriteResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
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

	err = h.voteServie.Unvote(ctx, user.UserID, uint(movieID))
	if err != nil {
		middleware.WriteResponse(w, http.StatusInternalServerError, "Failed to remove vote", nil)
		return
	}

	middleware.WriteResponse(w, http.StatusOK, "Vote removed successfully", nil)
}

func (h *VoteHandler) GetUserVotedMovies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		middleware.WriteResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	if !h.authService.IsLogin(ctx, user.UserID) {
		middleware.WriteResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	movies, err := h.voteServie.GetUserVotedMovies(ctx, user.UserID)
	if err != nil {
		middleware.WriteResponse(w, http.StatusInternalServerError, "Failed to fetch voted movies", nil)
		return
	}

	if movies == nil {
		movies = []entities.Movie{}
	}

	middleware.WriteResponse(w, http.StatusOK, "Voted movies fetched successfully", movies)
}
