package main

import (
	"flag"
	"log"
	"movie-festival-be/internal/app/entities"
	"movie-festival-be/internal/config"
	"movie-festival-be/internal/database"
	"movie-festival-be/package/logging"
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

}
