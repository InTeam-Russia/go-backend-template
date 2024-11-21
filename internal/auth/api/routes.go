package api

import (
	"github.com/InTeam-Russia/go-backend-template/internal/auth/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRoutes(r *gin.Engine, repo user.UserRepository, logger *zap.Logger) {
}
