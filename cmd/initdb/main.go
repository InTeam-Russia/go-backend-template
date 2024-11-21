package main

import (
	"os"

	"github.com/InTeam-Russia/go-backend-template/internal/auth"
	"github.com/InTeam-Russia/go-backend-template/internal/config"
	"github.com/InTeam-Russia/go-backend-template/internal/db"
	"github.com/InTeam-Russia/go-backend-template/internal/helpers"
)

func main() {
	config, err := config.LoadConfigFromEnv()
	logger := helpers.CreateLogger(config.LogLevel)

	pgPool, err := db.InitDb(config.PostgresUrl, os.Getenv("SQL_FILE"), logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer pgPool.Close()

	userRepo := auth.NewPgUserRepository(pgPool, logger)
	_, err = userRepo.Create(&auth.CreateUser{
		FirstName: "Admin",
		LastName:  "Admin",
		Username:  config.AdminUsername,
		Role:      "ADMIN",
		Password:  config.AdminPassword,
	})

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("Admin created!")
}
