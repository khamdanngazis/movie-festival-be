package main

import (
	"flag"
	"log"
	"movie-festival-be/internal/app/entities"
	"movie-festival-be/internal/config"
	"movie-festival-be/internal/database"
	"movie-festival-be/package/helper"
	"movie-festival-be/package/logging"
	"time"

	"gorm.io/gorm"
)

func main() {
	configFilePath := flag.String("config", "../config/config.yaml", "path to the config file")
	//logFile := flag.String("log.file", "../logs", "Logging file")

	flag.Parse()

	// Load the configuration
	cfg, err := config.LoadConfig(*configFilePath)
	if err != nil {
		logging.Log.Fatalf("Error loading configuration: %v", err)
	}
	logging.Log.Infof("Load configuration from %v", *configFilePath)

	db, err := database.InitDBPostgre(&cfg.Database.Main)

	if err != nil {
		logging.Log.Fatalf("Error initiate database connection: %v", err)
	}

	// Perform migration
	err = db.AutoMigrate(
		&entities.Movie{},
		&entities.User{},
		&entities.Vote{},
		&entities.Viewership{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migrated successfully!")
	insertAdminUser(db)
	insertMovies(db)

}

func insertAdminUser(db *gorm.DB) {
	// Check if an admin user already exists
	var user entities.User
	if err := db.Where("role = ?", "admin").First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No admin found, so create one
			hashPassword, _ := helper.HashPassword("Symantec2121")
			adminUser := entities.User{
				Name:      "Admin User",
				Email:     "admin@movie-festival.com",
				Password:  hashPassword, // Use hashed password here in real scenarios
				Role:      "admin",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := db.Create(&adminUser).Error; err != nil {
				logging.Log.Fatalf("Failed to insert admin user: %v", err)
			}

			log.Println("Admin user inserted successfully!")
		} else {
			logging.Log.Fatalf("Error checking for admin user: %v", err)
		}
	} else {
		log.Println("Admin user already exists.")
	}
}

func insertMovies(db *gorm.DB) {
	sampleMovies := []entities.Movie{
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

	for _, movie := range sampleMovies {
		if err := db.Create(&movie).Error; err != nil {
			logging.Log.Fatalf("Failed to insert movie '%s': %v", movie.Title, err)
		}
	}
}
