package handlers_test

import (
	"encoding/json"
	"movie-festival-be/internal/app/entities"
	"movie-festival-be/package/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SearchMovieResponse struct {
	middleware.HTTPResponse
	Data struct {
		Movies []entities.Movie `json:"movies"`
		Total  int              `json:"total"`
		Page   int              `json:"page"`
		Limit  int              `json:"limit"`
	} `json:"data"`
}

func TestSearchMovieHandler(t *testing.T) {

	httpRouter.GET("/movies/search", movieHandler.SearchMovies)

	req, err := http.NewRequest("GET", "/movies/search?query=Inception&page=1&limit=5", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	httpRouter.Mux().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response SearchMovieResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, response.Data.Total, 1)
	assert.Equal(t, len(response.Data.Movies), 1)
	assert.Equal(t, response.Data.Limit, 5)
	assert.Equal(t, response.Data.Page, 1)

}
