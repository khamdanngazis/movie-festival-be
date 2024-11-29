package handlers_test

import (
	"bytes"
	"encoding/json"
	"movie-festival-be/internal/app/entities"
	"movie-festival-be/internal/app/models"
	"movie-festival-be/package/middleware"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type movieResponse struct {
	middleware.HTTPResponse
	Movie entities.Movie `json:"data"`
}

func TestUpdateMovieHandler(t *testing.T) {

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

	request := `{
		"title": "Inception A",
		"description": "A mind-bending thriller by Christopher Nolan.a",
		"duration": 149,
		"artists": "Leonardo DiCaprio, Joseph Gordon-Levitt a",
		"genres": "Sci-Fi, Thriller a",
		"watch_url": "https://example.com/inceptiona"
	}`

	httpRouter.PUT("/movie/{id}", movieHandler.CreateMovie)

	req, err := http.NewRequest("PUT", "/movie/"+strconv.Itoa(int(movie.ID)), bytes.NewBuffer([]byte(request)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	httpRouter.Mux().Use(middleware.LoggingMiddleware)
	httpRouter.Mux().Use(middleware.AuthMiddleware)
	httpRouter.Mux().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response movieResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, response.Movie.ID)

}
