package repositories

import (
	"movie-festival-be/internal/app/entities"

	"gorm.io/gorm"
)

type MovieRepository interface {
	CreateMovie(movie *entities.Movie) error
	FindMovieByID(id uint) (*entities.Movie, error)
	UpdateMovie(movie *entities.Movie) error
	FindMovies(page, limit int) ([]entities.Movie, int64, error)
	SearchMovies(query string, limit int, offset int) ([]entities.Movie, int64, error)
}

type MovieRepositoryImpl struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &MovieRepositoryImpl{db: db}
}

func (r *MovieRepositoryImpl) CreateMovie(movie *entities.Movie) error {
	return r.db.Create(movie).Error
}

func (r *MovieRepositoryImpl) FindMovieByID(id uint) (*entities.Movie, error) {
	var movie entities.Movie
	err := r.db.First(&movie, id).Error
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *MovieRepositoryImpl) UpdateMovie(movie *entities.Movie) error {
	return r.db.Save(movie).Error
}

func (r *MovieRepositoryImpl) FindMovies(page, limit int) ([]entities.Movie, int64, error) {
	var movies []entities.Movie
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&entities.Movie{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&movies).Error
	if err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}

func (r *MovieRepositoryImpl) SearchMovies(query string, limit int, offset int) ([]entities.Movie, int64, error) {
	var movies []entities.Movie
	var total int64

	searchQuery := "%" + query + "%"

	// Count total matching records
	err := r.db.Model(&entities.Movie{}).
		Where("title LIKE ? OR description LIKE ? OR artists LIKE ? OR genres LIKE ?", searchQuery, searchQuery, searchQuery, searchQuery).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	err = r.db.Where("title LIKE ? OR description LIKE ? OR artists LIKE ? OR genres LIKE ?", searchQuery, searchQuery, searchQuery, searchQuery).
		Limit(limit).Offset(offset).
		Find(&movies).Error
	if err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}
