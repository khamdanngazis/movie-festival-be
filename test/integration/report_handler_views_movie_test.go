package handlers_test

import (
	"encoding/json"
	"movie-festival-be/internal/app/models"
	"movie-festival-be/package/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ReportResponse struct {
	middleware.HTTPResponse
	Report models.ReportViewsResponse `json:"data"`
}

func TestGetMovieViewsHandler(t *testing.T) {

	requestLogin := models.LoginRequest{
		Email:    "admin@movie-festival.com",
		Password: "Symantec2121",
	}

	token, _ := authService.Login(ctx, requestLogin)

	httpRouter.GET("/reports/views", reportHandler.GetReportViews)

	req, err := http.NewRequest("GET", "/reports/views", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	httpRouter.Mux().Use(middleware.AuthMiddleware)
	httpRouter.Mux().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response ReportResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	viewsMovie := 0
	mostViewMovie := ""
	for _, v := range sampleMovies {
		if v.Views > viewsMovie {
			viewsMovie = v.Views
			mostViewMovie = v.Title
		}
	}

	assert.Equal(t, mostViewMovie, response.Report.MostViewedMovie.Title)

}
