package services

import (
	"context"
	"errors"
	"movie-festival-be/internal/app/entities"
	"movie-festival-be/internal/app/models"
	"movie-festival-be/internal/app/repositories"
	"movie-festival-be/package/logging"
)

type MovieService interface {
	CreateMovie(ctx context.Context, request models.SaveMovieRequest) (*entities.Movie, error)
	UpdateMovie(ctx context.Context, id uint, req models.SaveMovieRequest) (*entities.Movie, error)
	ListMovies(ctx context.Context, page, limit int) ([]entities.Movie, int64, error)
	SearchMovies(ctx context.Context, query string, page int, limit int) ([]entities.Movie, int64, error)
	TrackView(ctx context.Context, movieID uint, ipAddress string, userID *uint) error
}

type MovieServiceImpl struct {
	movieRepo repositories.MovieRepository
}

func NewMovieService(movieRepo repositories.MovieRepository) MovieService {
	return &MovieServiceImpl{movieRepo: movieRepo}
}

func (s *MovieServiceImpl) CreateMovie(ctx context.Context, request models.SaveMovieRequest) (*entities.Movie, error) {
	movie := &entities.Movie{
		Title:       request.Title,
		Description: request.Description,
		Duration:    request.Duration,
		Artists:     request.Artists,
		Genres:      request.Genres,
		WatchURL:    request.WatchURL,
	}

	if err := s.movieRepo.CreateMovie(movie); err != nil {
		logging.LogError(ctx, "Error saving movie: %v", err)
		return nil, errors.New("internal server error")
	}

	return movie, nil
}

func (s *MovieServiceImpl) UpdateMovie(ctx context.Context, id uint, req models.SaveMovieRequest) (*entities.Movie, error) {
	movie, err := s.movieRepo.FindMovieByID(id)
	if err != nil {
		logging.LogError(ctx, "movie not found: %v", err)
		return nil, errors.New("movie not found")
	}

	movie.Title = req.Title
	movie.Description = req.Description
	movie.Duration = req.Duration
	movie.Artists = req.Artists
	movie.Genres = req.Genres
	movie.WatchURL = req.WatchURL

	if err := s.movieRepo.UpdateMovie(movie); err != nil {
		logging.LogError(ctx, "Error updating movie: %v", err)
		return nil, errors.New("internal server error")
	}

	return movie, nil
}

func (s *MovieServiceImpl) ListMovies(ctx context.Context, page, limit int) ([]entities.Movie, int64, error) {
	movies, total, err := s.movieRepo.FindMovies(page, limit)
	if err != nil {
		logging.LogError(ctx, "Error fetching movies: %v", err)
		return nil, 0, errors.New("failed to fetch movies")
	}
	return movies, total, nil
}

func (s *MovieServiceImpl) SearchMovies(ctx context.Context, query string, page int, limit int) ([]entities.Movie, int64, error) {
	offset := (page - 1) * limit
	movies, total, err := s.movieRepo.SearchMovies(query, limit, offset)
	if err != nil {
		logging.LogError(ctx, "Error searching movies: %v", err)
		return nil, 0, errors.New("internal server error")
	}
	return movies, total, nil
}

func (s *MovieServiceImpl) TrackView(ctx context.Context, movieID uint, ipAddress string, userID *uint) error {

	err := s.movieRepo.TrackViewership(movieID, ipAddress, userID)
	if err != nil {
		logging.LogError(ctx, "Error tracking viewership in transaction: %v", err)
		return errors.New("failed to track viewership")
	}

	return nil
}
