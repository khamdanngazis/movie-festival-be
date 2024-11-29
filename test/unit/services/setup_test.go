package services_test

import (
	"context"
	"fmt"
	"movie-festival-be/internal/app/entities"
	"movie-festival-be/internal/app/repositories"
	"movie-festival-be/internal/app/services"
	"movie-festival-be/internal/config"
	"movie-festival-be/internal/database"
	"movie-festival-be/package/helper"
	"movie-festival-be/package/logging"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	authRepo    repositories.AuthRepository
	authService services.AuthService

	movieRepo    repositories.MovieRepository
	movieService services.MovieService

	reportRepo    repositories.ReportRepository
	reportService services.ReportService

	voteRepo    repositories.VoteRepository
	voteService services.VoteService

	statsRepo    repositories.StatsRepository
	statsService services.StatsService

	ctx context.Context

	// Define sample movies as a global variable
	sampleMovies = []entities.Movie{
		{Title: "Inception", Description: "A mind-bending thriller by Christopher Nolan.", Duration: 148, Artists: "Leonardo DiCaprio, Joseph Gordon-Levitt", Genres: "Sci-Fi", WatchURL: "https://example.com/inception", Views: 5000, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "Interstellar", Description: "A science fiction epic exploring space and time.", Duration: 169, Artists: "Matthew McConaughey, Anne Hathaway", Genres: "Sci-Fi, Drama", WatchURL: "https://example.com/interstellar", Views: 4000, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "The Dark Knight", Description: "Batman faces his most challenging nemesis, The Joker.", Duration: 152, Artists: "Christian Bale, Heath Ledger", Genres: "Action, Thriller", WatchURL: "https://example.com/dark-knight", Views: 3000, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "The Matrix", Description: "A computer hacker discovers the true nature of reality.", Duration: 136, Artists: "Keanu Reeves, Laurence Fishburne", Genres: "Sci-Fi, Action", WatchURL: "https://example.com/matrix", Views: 4500, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "Titanic", Description: "A timeless love story set against the ill-fated maiden voyage.", Duration: 195, Artists: "Leonardo DiCaprio, Kate Winslet", Genres: "Drama, Romance", WatchURL: "https://example.com/titanic", Views: 6000, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "Avatar", Description: "A visually stunning epic on an alien planet.", Duration: 162, Artists: "Sam Worthington, Zoe Saldana", Genres: "Sci-Fi, Adventure", WatchURL: "https://example.com/avatar", Views: 5500, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "The Shawshank Redemption", Description: "The story of hope and friendship in prison.", Duration: 142, Artists: "Tim Robbins, Morgan Freeman", Genres: "Drama", WatchURL: "https://example.com/shawshank", Views: 7000, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "Pulp Fiction", Description: "An eclectic mix of stories woven into a masterpiece.", Duration: 154, Artists: "John Travolta, Uma Thurman", Genres: "Crime, Drama", WatchURL: "https://example.com/pulpfiction", Views: 3500, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "The Lord of the Rings: The Fellowship of the Ring", Description: "The epic journey begins to destroy the One Ring.", Duration: 178, Artists: "Elijah Wood, Ian McKellen", Genres: "Fantasy, Adventure", WatchURL: "https://example.com/lotr-fellowship", Views: 4800, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Title: "Forrest Gump", Description: "Life through the eyes of a simple, kind-hearted man.", Duration: 142, Artists: "Tom Hanks, Robin Wright", Genres: "Drama, Romance", WatchURL: "https://example.com/forrestgump", Views: 5200, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func initConfig() *config.Config {
	cfg, err := config.LoadConfig("../../../cmd/config/config-test.yaml")
	if err != nil {
		panic(err)
	}

	return cfg
}

func setup() {
	cfg := initConfig()
	db, err := database.InitDBPostgre(&cfg.Database.Main)
	if err != nil {
		panic(err)
	}

	authRepo = repositories.NewAuthRepository(db)
	authService = services.NewAuthService(authRepo)

	movieRepo = repositories.NewMovieRepository(db)
	movieService = services.NewMovieService(movieRepo)

	reportRepo = repositories.NewReportRepository(db)
	reportService = services.NewReportService(reportRepo)

	voteRepo = repositories.NewVoteRepository(db)
	voteService = services.NewVoteService(voteRepo)

	statsRepo = repositories.NewStatsRepository(db)
	statsService = services.NewStatsService(statsRepo)

	requestID := uuid.New().String()
	ctx = context.WithValue(context.Background(), logging.RequestIDKey, requestID)

	initData(cfg)
}

func initData(cfg *config.Config) {
	db, err := database.InitDBPostgre(&cfg.Database.Main)
	if err != nil {
		panic(err)
	}

	// Delete existing data
	deleteQueries := []string{
		"DELETE FROM users",
		"DELETE FROM movies",
		"DELETE FROM viewerships",
		"DELETE FROM votes",
	}
	DB, _ := db.DB()
	for _, query := range deleteQueries {
		_, err := DB.Exec(query)
		if err != nil {
			panic(fmt.Sprintf("Failed to execute query '%s': %v", query, err))
		}
	}

	// Insert data into users table
	hashPassword, _ := helper.HashPassword("Symantec2121")
	adminUser := entities.User{
		Name:      "Admin User",
		Email:     "admin@movie-festival.com",
		Password:  hashPassword,
		Role:      "admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&adminUser).Error; err != nil {
		logging.Log.Fatalf("Failed to insert admin user: %v", err)
	}

	for _, movie := range sampleMovies {
		if err := db.Create(&movie).Error; err != nil {
			logging.Log.Fatalf("Failed to insert movie '%s': %v", movie.Title, err)
		}
	}

	logging.Log.Infof("Sample data initialization completed successfully")
}
