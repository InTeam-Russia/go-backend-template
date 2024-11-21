package helpers

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func UUIDFromString(input string) (uuid.UUID, error) {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	hash := hasher.Sum(nil)

	var uuidBytes [16]byte
	copy(uuidBytes[:], hash[:16])

	uuidBytes[6] = (uuidBytes[6] & 0x0F) | 0x40
	uuidBytes[8] = (uuidBytes[8] & 0x3F) | 0x80

	return uuid.UUID(uuidBytes), nil
}

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
