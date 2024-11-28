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

	ctx context.Context
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

	logging.Log.Infof("Sample data initialization completed successfully")
}
