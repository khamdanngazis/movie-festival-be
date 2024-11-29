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
	"movie-festival-be/package/middleware"
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

	movieRepo := repositories.NewMovieRepository(db)
	movieService := services.NewMovieService(movieRepo)
	movieHandler := handlers.NewMoviesHandler(movieService, authService)

	httpRouter := router.NewMuxRouter()

	//ping handlers
	httpRouter.GET("/api/v1/ping", pinghandlers.Ping)

	//auth handler
	httpRouter.POST("/auth/login", authHandler.Login)

	//admin movie handler
	httpRouter.POSTWithMiddleware("/admin/movie", movieHandler.CreateMovie, middleware.AuthMiddleware)
	httpRouter.PUTWithMiddleware("/admin/movie/{id}", movieHandler.UpdateMovie, middleware.AuthMiddleware)

	//movie
	httpRouter.GET("/movies", movieHandler.ListMovies)
	httpRouter.GET("/movies/search", movieHandler.SearchMovies)

	httpRouter.SERVE(cfg.AppPort)
}
func initLogging() {
	logging.InitLogger()
	logging.Log.SetOutput(os.Stdout)
}
