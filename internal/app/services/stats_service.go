package services

import (
	"context"
	"movie-festival-be/internal/app/models"
	"movie-festival-be/internal/app/repositories"
	"movie-festival-be/package/logging"
)

type StatsService interface {
	GetAdminStats(ctx context.Context) (*models.AdminStatsResponse, error)
}

type StatsServiceImpl struct {
	statsRepo repositories.StatsRepository
}

func NewStatsService(statsRepo repositories.StatsRepository) StatsService {
	return &StatsServiceImpl{statsRepo: statsRepo}
}

func (s *StatsServiceImpl) GetAdminStats(ctx context.Context) (*models.AdminStatsResponse, error) {
	// Fetch most voted movie
	movie, votesCount, err := s.statsRepo.GetMostVotedMovie()
	if err != nil {
		logging.LogError(ctx, "Error fetching most voted movie: %v", err)
		return nil, err
	}

	// Fetch most viewed genre
	genre, totalViews, err := s.statsRepo.GetMostViewedGenre()
	if err != nil {
		logging.LogError(ctx, "Error fetching most viewed genre: %v", err)
		return nil, err
	}

	response := &models.AdminStatsResponse{
		MostVotedMovie: models.MovieStatsResponse{
			ID:          movie.ID,
			Title:       movie.Title,
			Description: movie.Description,
			Duration:    movie.Duration,
			Artists:     movie.Artists,
			Genres:      movie.Genres,
			VotesCount:  votesCount,
		},
		MostViewedGenre: models.GenreStatsResponse{
			Genre:      genre,
			TotalViews: totalViews,
		},
	}

	return response, nil
}
