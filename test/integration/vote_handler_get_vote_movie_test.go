package handlers_test

import (
	"encoding/json"
	"movie-festival-be/internal/app/entities"
	"movie-festival-be/internal/app/models"
	"movie-festival-be/package/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ViewVoteResponse struct {
	middleware.HTTPResponse
	Data []entities.Movie `json:"data"`
}

func TestGetMovieHandler(t *testing.T) {

	user, _ := authRepo.FindUserByEmail("admin@movie-festival.com")

	request := models.SaveMovieRequest{
		Title:       "Inception",
		Description: "A mind-bending thriller by Christopher Nolan.",
		Duration:    148,
		Artists:     "Leonardo DiCaprio, Joseph Gordon-Levitt",
		Genres:      "Sci-Fi, Thriller",
		WatchURL:    "https://example.com/inception",
	}

	movie, _ := movieService.CreateMovie(ctx, request)

	voteService.VoteMovie(ctx, movie.ID, user.ID)

	requestLogin := models.LoginRequest{
		Email:    "admin@movie-festival.com",
		Password: "Symantec2121",
	}

	token, _ := authService.Login(ctx, requestLogin)

	httpRouter.GETWithMiddleware("/users/me/votes", voteHandler.GetUserVotedMovies, middleware.AuthMiddleware)

	req, err := http.NewRequest("GET", "/users/me/votes", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	httpRouter.Mux().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response ViewVoteResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(response.Data), 1)
	assert.Equal(t, movie.ID, response.Data[0].ID)

}
