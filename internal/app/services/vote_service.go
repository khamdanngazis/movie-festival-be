package services

import (
	"context"
	"errors"
	"movie-festival-be/internal/app/entities"
	"movie-festival-be/internal/app/repositories"
	"movie-festival-be/package/logging"
)

type VoteService interface {
	VoteMovie(ctx context.Context, movieID, userID uint) error
	Unvote(ctx context.Context, userID, movieID uint) error
	GetUserVotedMovies(ctx context.Context, userID uint) ([]entities.Movie, error)
}

type VoteServiceImpl struct {
	voteRepo repositories.VoteRepository
}

func NewVoteService(voteRepo repositories.VoteRepository) VoteService {
	return &VoteServiceImpl{voteRepo: voteRepo}
}

func (s *VoteServiceImpl) VoteMovie(ctx context.Context, movieID, userID uint) error {
	// Check if the user has already voted
	hasVoted, err := s.voteRepo.HasUserVoted(movieID, userID)
	if err != nil {
		logging.LogError(ctx, "Error checking vote status: %v", err)
		return errors.New("internal server error")
	}

	if hasVoted {
		return errors.New("user has already voted for this movie")
	}

	// Create a new vote
	vote := &entities.Vote{
		UserID:  userID,
		MovieID: movieID,
	}
	if err := s.voteRepo.CreateVote(vote); err != nil {
		logging.LogError(ctx, "Error saving vote: %v", err)
		return errors.New("failed to cast vote")
	}

	return nil
}

func (s *VoteServiceImpl) Unvote(ctx context.Context, userID, movieID uint) error {
	err := s.voteRepo.Unvote(userID, movieID)
	if err != nil {
		return errors.New("failed to remove vote")
	}
	return nil
}

func (s *VoteServiceImpl) GetUserVotedMovies(ctx context.Context, userID uint) ([]entities.Movie, error) {
	movies, err := s.voteRepo.GetUserVotes(userID)
	if err != nil {
		logging.LogError(ctx, "Error fetching user votes: %v", err)
		return nil, errors.New("failed to fetch voted movies")
	}
	return movies, nil
}
