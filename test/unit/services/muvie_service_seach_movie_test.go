package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchMovieService(t *testing.T) {
	list, total, err := movieService.SearchMovies(ctx, sampleMovies[0].Title, 1, 5)

	assert.Nil(t, err)
	assert.Equal(t, int(total), 1)
	assert.Equal(t, len(list), 1)
	assert.Equal(t, list[0].Title, sampleMovies[0].Title)

}
