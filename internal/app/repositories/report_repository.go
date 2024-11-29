package repositories

import (
	"movie-festival-be/internal/app/entities"

	"gorm.io/gorm"
)

type ReportRepository interface {
	GetMostViewedMovie() (*entities.Movie, error)
	GetGenreViewCounts() ([]GenreViewCount, error)
}

type ReportRepositoryImpl struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &ReportRepositoryImpl{db: db}
}

type GenreViewCount struct {
	Genre      string
	TotalViews int
}

func (r *ReportRepositoryImpl) GetMostViewedMovie() (*entities.Movie, error) {
	var movie entities.Movie
	err := r.db.Order("views DESC").First(&movie).Error
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *ReportRepositoryImpl) GetGenreViewCounts() ([]GenreViewCount, error) {
	var results []GenreViewCount
	err := r.db.
		Model(&entities.Movie{}).
		Select("genres AS genre, SUM(views) AS total_views").
		Group("genres").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
