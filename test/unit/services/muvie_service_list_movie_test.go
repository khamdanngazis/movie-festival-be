package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListMovieService(t *testing.T) {
	list, total, err := movieService.ListMovies(ctx, 1, 5)
	assert.Nil(t, err)
	assert.Equal(t, int(total), len(sampleMovies))
	assert.Equal(t, len(list), 5)

}
