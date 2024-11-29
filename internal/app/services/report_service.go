package services

import (
	"context"
	"errors"
	"movie-festival-be/internal/app/models"
	"movie-festival-be/internal/app/repositories"
)

type ReportService interface {
	GetReportViews(ctx context.Context) (*models.ReportViewsResponse, error)
}

type ReportServiceImpl struct {
	reportRepo repositories.ReportRepository
}

func NewReportService(reportRepo repositories.ReportRepository) ReportService {
	return &ReportServiceImpl{reportRepo: reportRepo}
}

func (s *ReportServiceImpl) GetReportViews(ctx context.Context) (*models.ReportViewsResponse, error) {
	movie, err := s.reportRepo.GetMostViewedMovie()
	if err != nil {
		return nil, errors.New("could not fetch most viewed movie")
	}

	genreStats, err := s.reportRepo.GetGenreViewCounts()
	if err != nil {
		return nil, errors.New("could not fetch genre statistics")
	}

	// Transform the data
	var genreStatsResponse []models.GenreStatItem
	for _, genreStat := range genreStats {
		genreStatsResponse = append(genreStatsResponse, models.GenreStatItem{
			Genre:      genreStat.Genre,
			TotalViews: genreStat.TotalViews,
		})
	}

	return &models.ReportViewsResponse{
		MostViewedMovie: models.MovieViewResponse{
			ID:          movie.ID,
			Title:       movie.Title,
			Artists:     movie.Artists,
			Genres:      movie.Genres,
			Description: movie.Description,
			Duration:    movie.Duration,
			ViewCount:   int64(movie.Views),
		},
		GenreStats: genreStatsResponse,
	}, nil
}
