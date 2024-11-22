package main

import (
	"os"
	"time"

	"github.com/InTeam-Russia/go-backend-template/internal/auth"
	"github.com/InTeam-Russia/go-backend-template/internal/config"
	"github.com/InTeam-Russia/go-backend-template/internal/db"
	"github.com/InTeam-Russia/go-backend-template/internal/helpers"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	config, err := config.LoadConfigFromEnv()
	if err != nil {
		panic(err)
	}

	logger := helpers.CreateLogger(config.LogLevel)

	pgPool, err := db.CreatePool(config.PostgresUrl, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer pgPool.Close()

	redisOpts, err := redis.ParseURL(config.RedisUrl)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	redisClient := redis.NewClient(redisOpts)
	defer redisClient.Close()

	r := gin.New()
	helpers.SetupCORS(r, config)
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	cookieConfig := auth.DefaultCookieConfig()
	cookieConfig.Secure = config.SessionCookieSecure
	cookieConfig.Domain = config.SessionCookieDomain

	userRepo := auth.NewPgUserRepository(pgPool, logger)
	sessionRepo := auth.NewRedisSessionRepository(redisClient, logger)

	auth.SetupRoutes(r, userRepo, sessionRepo, logger, cookieConfig)

	r.Run()
}
