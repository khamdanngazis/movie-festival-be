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
