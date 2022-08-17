package main

import (
	"log"

	"github.com/samthehai/ml-backend-test-samthehai/api"
	"github.com/samthehai/ml-backend-test-samthehai/config"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/db/mysql"
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

	connManager, closeAllConn, err := mysql.NewConnManager(cfg)
	if err != nil {
		log.Fatalf("failed to init mysql client: %v", err)
	}
	defer closeAllConn()

	s := api.NewServer(cfg, appLogger, connManager)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
