package main

import (
	"log"

	"github.com/samthehai/ml-backend-test-samthehai/api"
	"github.com/samthehai/ml-backend-test-samthehai/config"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/db/mysql"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/logger"
	"github.com/spf13/viper"
)

const defaultConfigPath = "config"

// @title MonstarLab Backend Test API
// @version 1.0
// @description MonstarLab Backend Test REST API written in Golang
//
// @contact.name Hai Sam
// @contact.url https://github.com/samthehai
// @contact.email samthehai@gmail.com
//
// @host localhost:5000
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	configPath := viper.GetString("CONFIG_PATH")
	if len(configPath) == 0 {
		configPath = defaultConfigPath
	}

	cfg, err := config.LoadConfig(configPath)
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
