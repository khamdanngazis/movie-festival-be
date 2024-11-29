package repositories

import (
	"movie-festival-be/internal/app/entities"

	"gorm.io/gorm"
)

type VoteRepository interface {
	HasUserVoted(movieID, userID uint) (bool, error)
	CreateVote(vote *entities.Vote) error
	Unvote(userID, movieID uint) error
	GetUserVotes(userID uint) ([]entities.Movie, error)
}

type VoteRepositoryImpl struct {
	db *gorm.DB
}

func NewVoteRepository(db *gorm.DB) VoteRepository {
	return &VoteRepositoryImpl{db: db}
}

func (r *VoteRepositoryImpl) HasUserVoted(movieID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entities.Vote{}).
		Where("movie_id = ? AND user_id = ?", movieID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *VoteRepositoryImpl) CreateVote(vote *entities.Vote) error {
	return r.db.Create(vote).Error
}

func (r *VoteRepositoryImpl) Unvote(userID, movieID uint) error {
	return r.db.Unscoped().Where("user_id = ? AND movie_id = ?", userID, movieID).Delete(&entities.Vote{}).Error
}

func (r *VoteRepositoryImpl) GetUserVotes(userID uint) ([]entities.Movie, error) {
	var movies []entities.Movie

	// Join votes with movies to fetch voted movies
	err := r.db.Table("movies").
		Select("movies.*").
		Joins("join votes on votes.movie_id = movies.id").
		Where("votes.user_id = ?", userID).
		Scan(&movies).Error

	return movies, err
}
