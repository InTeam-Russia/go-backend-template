package main

import (
	"os"

	"github.com/InTeam-Russia/go-backend-template/internal/config"
	"github.com/InTeam-Russia/go-backend-template/internal/db"
	"github.com/InTeam-Russia/go-backend-template/internal/helpers"
)

func main() {
	config, err := config.LoadConfigFromEnv()
	logger := helpers.CreateLogger(config.LogLevel)

	pgPool, err := db.DropDb(config.PostgresUrl, os.Getenv("SQL_FILE"), logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer pgPool.Close()
}
