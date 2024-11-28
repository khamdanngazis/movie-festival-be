package database

import (
	"fmt"
	"movie-festival-be/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDBPostgre(cfg *config.Postgre) (*gorm.DB, error) {
	config := &gorm.Config{
		PrepareStmt: true,
	}
	if cfg.Debug {
		config.Logger = logger.Default.LogMode(logger.Info) // Set logger level to Info
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Host, cfg.User, cfg.DBName, cfg.Password, cfg.Port, cfg.SSLMode, cfg.TimeZone,
	)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), config)

	if err != nil {
		return nil, err
	}

	// Configure connection pooling
	sqlDB, err := db.DB() // Get the underlying sql.DB from GORM
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)                 // Maximum open connections
	sqlDB.SetMaxIdleConns(10)                 // Maximum idle connections
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of connection

	return db, nil
}
