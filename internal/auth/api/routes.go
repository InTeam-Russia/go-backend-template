package api

import (
	"net/http"

	"github.com/InTeam-Russia/go-backend-template/internal/apierr"
	"github.com/InTeam-Russia/go-backend-template/internal/auth/session"
	"github.com/InTeam-Russia/go-backend-template/internal/auth/shared"
	"github.com/InTeam-Russia/go-backend-template/internal/auth/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Seconds = int

type CookieConfig struct {
	SessionLifetime Seconds
	Secure          bool
	Path            string
	HttpOnly        bool
	Domain          string
}

func DefaultCookieConfig() *CookieConfig {
	return &CookieConfig{
		SessionLifetime: 604800, // 7 weeks
		Secure:          true,
		Path:            "/",
		HttpOnly:        true,
		Domain:          "",
	}
}

func SetupRoutes(
	r *gin.Engine,
	userRepo user.UserRepository,
	sessionRepo session.SessionRepository,
	logger *zap.Logger,
	cookieConfig *CookieConfig,
) {
	r.POST("/login", func(c *gin.Context) {
		var json Login
		if err := c.ShouldBindBodyWithJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, apierr.InvalidJsonError)
			return
		}

		user, err := userRepo.GetByUsername(json.Username)

		if err != nil {
			c.JSON(http.StatusInternalServerError, apierr.InternalServerError)
			logger.Error(err.Error())
			return
		}

		if user == nil {
			c.JSON(http.StatusNotFound, apierr.NotFoundError)
			return
		}

		if !shared.ValidPassword(json.Password, user.PasswordHash, user.PasswordSalt) {
			c.JSON(http.StatusUnauthorized, apierr.WrongCredentials)
			return
		}

		session, err := sessionRepo.Create(user.Id, cookieConfig.SessionLifetime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, apierr.InternalServerError)
			logger.Error(err.Error())
			return
		}

		c.SetCookie(
			"SESSION_ID",
			session.Id.String(),
			cookieConfig.SessionLifetime,
			cookieConfig.Path,
			cookieConfig.Domain,
			cookieConfig.Secure,
			cookieConfig.HttpOnly,
		)

		c.JSON(http.StatusCreated, gin.H{
			"status": "OK",
		})
	})

	r.POST("/register", func(c *gin.Context) {
		var json Register
		if err := c.ShouldBindBodyWithJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, apierr.InvalidJsonError)
			return
		}

		createUser := user.CreateUser{
			FirstName: json.FirstName,
			LastName:  json.LastName,
			Username:  json.Username,
			Password:  json.Password,
			Role:      "USER",
		}

		u, err := userRepo.Create(&createUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, apierr.InternalServerError)
			logger.Error(err.Error())
			return
		}

		c.JSON(http.StatusCreated, user.UserOut{
			Id:        u.Id,
			CreatedAt: u.CreatedAt,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Username:  u.Username,
			Role:      u.Role,
		})
	})
}
