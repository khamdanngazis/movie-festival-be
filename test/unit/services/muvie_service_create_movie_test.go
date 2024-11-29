package services_test

import (
	"movie-festival-be/internal/app/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMovieService(t *testing.T) {
	request := models.SaveMovieRequest{
		Title:       "Inception",
		Description: "A mind-bending thriller by Christopher Nolan.",
		Duration:    148,
		Artists:     "Leonardo DiCaprio, Joseph Gordon-Levitt",
		Genres:      "Sci-Fi, Thriller",
		WatchURL:    "https://example.com/inception",
	}

	movie, err := movieService.CreateMovie(ctx, request)
	assert.Nil(t, err)
	assert.NotEmpty(t, movie.ID)
	assert.Equal(t, request.Title, movie.Title)
	assert.Equal(t, request.Description, movie.Description)
	assert.Equal(t, request.Duration, movie.Duration)
	assert.Equal(t, request.Artists, movie.Artists)
	assert.Equal(t, request.Genres, movie.Genres)
	assert.Equal(t, request.WatchURL, movie.WatchURL)

}
