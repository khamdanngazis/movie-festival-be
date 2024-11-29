package services_test

import (
	"movie-festival-be/internal/app/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateMovieService(t *testing.T) {
	request := models.SaveMovieRequest{
		Title:       "Inception",
		Description: "A mind-bending thriller by Christopher Nolan.",
		Duration:    148,
		Artists:     "Leonardo DiCaprio, Joseph Gordon-Levitt",
		Genres:      "Sci-Fi, Thriller",
		WatchURL:    "https://example.com/inception",
	}

	movie, _ := movieService.CreateMovie(ctx, request)

	requestUpdate := models.SaveMovieRequest{
		Title:       "Inception A",
		Description: "A mind-bending thriller by Christopher Nolan. A",
		Duration:    144,
		Artists:     "Leonardo DiCaprio, Joseph Gordon-Levitt t",
		Genres:      "Sci-Fi, Thriller t",
		WatchURL:    "https://example.com/inceptiona",
	}

	movieUpdate, err := movieService.UpdateMovie(ctx, movie.ID, requestUpdate)
	assert.Nil(t, err)
	assert.Equal(t, movie.ID, movieUpdate.ID)
	assert.Equal(t, requestUpdate.Title, movieUpdate.Title)
	assert.Equal(t, requestUpdate.Description, movieUpdate.Description)
	assert.Equal(t, requestUpdate.Duration, movieUpdate.Duration)
	assert.Equal(t, requestUpdate.Artists, movieUpdate.Artists)
	assert.Equal(t, requestUpdate.Genres, movieUpdate.Genres)
	assert.Equal(t, requestUpdate.WatchURL, movieUpdate.WatchURL)

}
