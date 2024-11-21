package helpers

import (
	"crypto/sha1"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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