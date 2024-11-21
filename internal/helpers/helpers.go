package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupCORS(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
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
