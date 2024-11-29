package handlers_test

import (
	"bytes"
	"encoding/json"
	"movie-festival-be/internal/app/entities"
	"movie-festival-be/internal/app/models"
	"movie-festival-be/package/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MovieResponse struct {
	middleware.HTTPResponse
	Movie entities.Movie `json:"data"`
}

func TestCreateMovieHandler(t *testing.T) {

	requestLogin := models.LoginRequest{
		Email:    "admin@movie-festival.com",
		Password: "Symantec2121",
	}

	token, _ := authService.Login(ctx, requestLogin)

	request := `{
		"title": "Inception",
		"description": "A mind-bending thriller by Christopher Nolan.",
		"duration": 148,
		"artists": "Leonardo DiCaprio, Joseph Gordon-Levitt",
		"genres": "Sci-Fi, Thriller",
		"watch_url": "https://example.com/inception"
	}`

	httpRouter.POST("/movie", movieHandler.CreateMovie)

	req, err := http.NewRequest("POST", "/movie", bytes.NewBuffer([]byte(request)))
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

	var response MovieResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, response.Movie.ID)

}
