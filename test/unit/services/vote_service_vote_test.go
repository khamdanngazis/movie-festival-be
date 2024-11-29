package services_test

import (
	"movie-festival-be/internal/app/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVoteMovieService(t *testing.T) {
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

	err := voteService.VoteMovie(ctx, movie.ID, user.ID)

	assert.Nil(t, err)
	hasUserVoted, _ := voteRepo.HasUserVoted(movie.ID, user.ID)
	assert.True(t, hasUserVoted)

}
