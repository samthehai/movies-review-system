package main

import (
	"log"

	"github.com/samthehai/ml-backend-test-samthehai/config"
	"github.com/samthehai/ml-backend-test-samthehai/internal/server"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof("LogLevel: %s, Mode: %s", cfg.Logger.Level, cfg.Server.Mode)
	s := server.NewServer(cfg, appLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
