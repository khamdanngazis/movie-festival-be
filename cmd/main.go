package main

import (
	"flag"
	"movie-festival-be/internal/app/repositories"
	"movie-festival-be/internal/app/services"
	"movie-festival-be/internal/config"
	"movie-festival-be/internal/database"
	"movie-festival-be/internal/interface/handlers"
	"movie-festival-be/internal/interface/router"
	"movie-festival-be/package/logging"
	"os"
)

func main() {
	configFilePath := flag.String("config", "config/config.yaml", "path to the config file")
	//logFile := flag.String("log.file", "../logs", "Logging file")

	flag.Parse()

	initLogging()

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

	pinghandlers := handlers.NewPinghandlers()

	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authHandler := handlers.NewAuthhandler(authService)

	httpRouter := router.NewMuxRouter()

	//ping handlers
	httpRouter.GET("/api/v1/ping", pinghandlers.Ping)

	//auth handler
	httpRouter.POST("/auth/login", authHandler.Login)

	httpRouter.SERVE(cfg.AppPort)
}
func initLogging() {
	logging.InitLogger()
	logging.Log.SetOutput(os.Stdout)
}
