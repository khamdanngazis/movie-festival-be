package handlers_test

import (
	"bytes"
	"encoding/json"
	"movie-festival-be/internal/app/models"
	"movie-festival-be/package/middleware"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrackViewMovieHandler(t *testing.T) {

	requestInsert := models.SaveMovieRequest{
		Title:       "Inception",
		Description: "A mind-bending thriller by Christopher Nolan.",
		Duration:    148,
		Artists:     "Leonardo DiCaprio, Joseph Gordon-Levitt",
		Genres:      "Sci-Fi, Thriller",
		WatchURL:    "https://example.com/inception",
	}

	movie, _ := movieService.CreateMovie(ctx, requestInsert)

	request := `{
		"watch_duration" : 60
	}`

	httpRouter.POSTWithMiddleware("/movies/{id:[0-9]+}/view", movieHandler.TrackView, middleware.GuestMiddleware)

	req, err := http.NewRequest("POST", "/movies/"+strconv.Itoa(int(movie.ID))+"/view", bytes.NewBuffer([]byte(request)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Forwarded-For", "127.0.0.1")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	httpRouter.Mux().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response middleware.HTTPResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}
}
