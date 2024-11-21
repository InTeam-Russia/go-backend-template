package main

import (
	"os"

	"github.com/InTeam-Russia/go-backend-template/internal/config"
	"github.com/InTeam-Russia/go-backend-template/internal/db"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	config, err := config.LoadConfigFromEnv(logger)

	pgPool, err := db.InitDb(config.PostgresUrl, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer pgPool.Close()
}
