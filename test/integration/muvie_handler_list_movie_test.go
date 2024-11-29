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

type ListMovieResponse struct {
	middleware.HTTPResponse
	Data struct {
		Movies []entities.Movie `json:"movies"`
		Total  int              `json:"total"`
		Page   int              `json:"page"`
		Limit  int              `json:"limit"`
	} `json:"data"`
}

func TestListMovieHandler(t *testing.T) {

	httpRouter.GET("/movies", movieHandler.ListMovies)

	req, err := http.NewRequest("GET", "/movies?page=1&limit=5", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	httpRouter.Mux().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response ListMovieResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, response.Data.Total, len(sampleMovies))
	assert.Equal(t, len(response.Data.Movies), 5)
	assert.Equal(t, response.Data.Limit, 5)
	assert.Equal(t, response.Data.Page, 1)

}
