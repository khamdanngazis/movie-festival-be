package services_test

import (
	"movie-festival-be/internal/app/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrackViewMovieService(t *testing.T) {
	request := models.SaveMovieRequest{
		Title:       "Inception",
		Description: "A mind-bending thriller by Christopher Nolan.",
		Duration:    148,
		Artists:     "Leonardo DiCaprio, Joseph Gordon-Levitt",
		Genres:      "Sci-Fi, Thriller",
		WatchURL:    "https://example.com/inception",
	}

	movie, _ := movieService.CreateMovie(ctx, request)
	ipAddress := "127.0.0.9"

	err := movieService.TrackView(ctx, uint(movie.ID), ipAddress, nil)

	assert.Nil(t, err)
	newRecordMovie, _ := movieRepo.FindMovieByID(movie.ID)
	assert.Equal(t, movie.Views+1, newRecordMovie.Views)
	viewership, _ := movieRepo.FindViewershipByMovieID(movie.ID)
	assert.Equal(t, viewership.IPAddress, ipAddress)

}
