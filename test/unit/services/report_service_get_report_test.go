package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetReportViewsService(t *testing.T) {
	// Call the service method
	views, err := reportService.GetReportViews(ctx)

	// Assert no error occurred
	assert.Nil(t, err)
	assert.NotNil(t, views)
	viewsMovie := 0
	mostViewMovie := ""
	for _, v := range sampleMovies {
		if v.Views > viewsMovie {
			viewsMovie = v.Views
			mostViewMovie = v.Title
		}
	}

	assert.Equal(t, mostViewMovie, views.MostViewedMovie.Title)
}
