package main

import (
	"flag"
	"movie-festival-be/internal/config"
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

	pinghandlers := handlers.NewPinghandlers()

	httpRouter := router.NewMuxRouter()

	//ping handlers
	httpRouter.GET("/api/v1/ping", pinghandlers.Ping)

	httpRouter.SERVE(cfg.AppPort)
}
func initLogging() {
	logging.InitLogger()
	logging.Log.SetOutput(os.Stdout)
}
