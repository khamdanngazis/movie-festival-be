package handlers_test

import (
	"encoding/json"
	"movie-festival-be/internal/app/models"
	"movie-festival-be/package/middleware"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVoteMovieHandler(t *testing.T) {

	requestInsert := models.SaveMovieRequest{
		Title:       "Inception",
		Description: "A mind-bending thriller by Christopher Nolan.",
		Duration:    148,
		Artists:     "Leonardo DiCaprio, Joseph Gordon-Levitt",
		Genres:      "Sci-Fi, Thriller",
		WatchURL:    "https://example.com/inception",
	}

	movie, _ := movieService.CreateMovie(ctx, requestInsert)

	requestLogin := models.LoginRequest{
		Email:    "admin@movie-festival.com",
		Password: "Symantec2121",
	}

	token, _ := authService.Login(ctx, requestLogin)

	httpRouter.POSTWithMiddleware("/movies/{id:[0-9]+}/vote", voteHandler.VoteMovie, middleware.AuthMiddleware)

	req, err := http.NewRequest("POST", "/movies/"+strconv.Itoa(int(movie.ID))+"/vote", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	httpRouter.Mux().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response middleware.HTTPResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	user, _ := authRepo.FindUserByEmail("admin@movie-festival.com")
	hasUserVoted, _ := voteRepo.HasUserVoted(movie.ID, user.ID)
	assert.True(t, hasUserVoted)

}
