package repositories

import (
	"errors"
	"log"
	"movie-festival-be/internal/app/entities"

	"gorm.io/gorm"
)

type MovieRepository interface {
	CreateMovie(movie *entities.Movie) error
	FindMovieByID(id uint) (*entities.Movie, error)
	UpdateMovie(movie *entities.Movie) error
	FindMovies(page, limit int) ([]entities.Movie, int64, error)
	SearchMovies(query string, limit int, offset int) ([]entities.Movie, int64, error)
	TrackViewership(movieID uint, ipAddress string, userID *uint, WatchDuration int) error
	FindViewershipByMovieID(movieID uint) (*entities.Viewership, error)
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

func (r *MovieRepositoryImpl) FindViewershipByMovieID(movieID uint) (*entities.Viewership, error) {
	var viewership entities.Viewership
	err := r.db.Where("movie_id = ?", movieID).First(&viewership).Error
	if err != nil {
		return nil, err
	}
	return &viewership, nil
}

func (r *MovieRepositoryImpl) TrackViewership(movieID uint, ipAddress string, userID *uint, watchDuration int) error {
	if watchDuration <= 0 {
		return errors.New("watch duration must be greater than zero")
	}

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err := tx.Error; err != nil {
			tx.Rollback()
		}
	}()

	if watchDuration >= 60 {
		if err := tx.Model(&entities.Movie{}).Where("id = ?", movieID).
			Update("views", gorm.Expr("views + 1")).Error; err != nil {
			log.Printf("Failed to increment views for movie ID %d: %v", movieID, err)
			tx.Rollback()
			return err
		}
	}

	viewership := entities.Viewership{
		MovieID:       movieID,
		IPAddress:     ipAddress,
		UserID:        userID,
		WatchDuration: watchDuration,
	}
	if err := tx.Create(&viewership).Error; err != nil {
		log.Printf("Failed to create viewership record: %v", err)
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
