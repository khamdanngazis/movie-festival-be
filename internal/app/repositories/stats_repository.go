package repositories

import (
	"movie-festival-be/internal/app/entities"

	"gorm.io/gorm"
)

type StatsRepository interface {
	GetMostVotedMovie() (*entities.Movie, int64, error)
	GetMostViewedGenre() (string, int64, error)
}

type StatsRepositoryImpl struct {
	db *gorm.DB
}

func NewStatsRepository(db *gorm.DB) StatsRepository {
	return &StatsRepositoryImpl{db: db}
}

func (r *StatsRepositoryImpl) GetMostVotedMovie() (*entities.Movie, int64, error) {
	var movie entities.Movie
	var votesCount int64

	err := r.db.Table("movies").
		Select("movies.*, COUNT(votes.id) as votes_count").
		Joins("JOIN votes ON votes.movie_id = movies.id").
		Group("movies.id").
		Order("votes_count DESC").
		Limit(1).
		Scan(&movie).Error

	if err != nil {
		return nil, 0, err
	}

	err = r.db.Table("votes").
		Where("movie_id = ?", movie.ID).
		Count(&votesCount).Error

	return &movie, votesCount, err
}

func (r *StatsRepositoryImpl) GetMostViewedGenre() (string, int64, error) {
	var result struct {
		Genre      string
		TotalViews int64
	}

	err := r.db.Table("movies").
		Select("genres as genre, SUM(views) as total_views").
		Group("genres").
		Order("total_views DESC").
		Limit(1).
		Scan(&result).Error

	return result.Genre, result.TotalViews, err
}
