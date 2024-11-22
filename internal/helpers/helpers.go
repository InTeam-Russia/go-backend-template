package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/InTeam-Russia/go-backend-template/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupCORS(r *gin.Engine, config *config.Config) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.AllowOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
}

func CreateLogger(logLevel string) *zap.Logger {
	rawJSON := []byte(fmt.Sprintf(
		`{
	   "level": "%s",
	   "encoding": "console",
	   "outputPaths": ["stdout"],
	   "errorOutputPaths": ["stderr"],
	   "encoderConfig": {
	     "messageKey": "message",
	     "levelKey": "level",
	     "levelEncoder": "lowercase"
	   }
	  }`,
		logLevel,
	))

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger
}
