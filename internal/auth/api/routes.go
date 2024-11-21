package api

import (
	"net/http"

	"github.com/InTeam-Russia/go-backend-template/internal/apierr"
	"github.com/InTeam-Russia/go-backend-template/internal/auth/session"
	"github.com/InTeam-Russia/go-backend-template/internal/auth/shared"
	"github.com/InTeam-Russia/go-backend-template/internal/auth/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const sessionCookieName = "SESSION_ID"

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
		var loginJson Login
		if err := c.ShouldBindBodyWithJSON(&loginJson); err != nil {
			c.JSON(http.StatusBadRequest, apierr.InvalidJsonError)
			return
		}

		user, err := userRepo.GetByUsername(loginJson.Username)

		if err != nil {
			c.JSON(http.StatusInternalServerError, apierr.InternalServerError)
			logger.Error(err.Error())
			return
		}

		if user == nil {
			c.JSON(http.StatusNotFound, apierr.NotFoundError)
			return
		}

		if !shared.ValidPassword(loginJson.Password, user.PasswordHash, user.PasswordSalt) {
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
			sessionCookieName,
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
		var registerJson Register
		if err := c.ShouldBindBodyWithJSON(&registerJson); err != nil {
			c.JSON(http.StatusBadRequest, apierr.InvalidJsonError)
			return
		}

		createUser := user.CreateUser{
			FirstName: registerJson.FirstName,
			LastName:  registerJson.LastName,
			Username:  registerJson.Username,
			Password:  registerJson.Password,
			Role:      "USER",
		}

		u, err := userRepo.Create(&createUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, apierr.InternalServerError)
			logger.Error(err.Error())
			return
		}

		c.JSON(http.StatusCreated, mapUserToUserOut(u))
	})

	r.POST("/logout", func(c *gin.Context) {
		cookie, err := c.Cookie(sessionCookieName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, apierr.CookieNotExists)
			return
		}

		cookieIdUUID, err := uuid.Parse(cookie)
		if err != nil {
			c.JSON(http.StatusUnauthorized, apierr.CookieNotExists)
			return
		}

		err = sessionRepo.DeleteById(cookieIdUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, apierr.InternalServerError)
			logger.Error(err.Error())
			return
		}

		c.SetCookie(sessionCookieName, "", -1, "/", "localhost", false, true)

		c.JSON(http.StatusCreated, gin.H{
			"status": "OK",
		})
	})

	r.GET("/session", func(c *gin.Context) {
		cookie, err := c.Cookie(sessionCookieName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, apierr.CookieNotExists)
			return
		}

		cookieIdUUID, err := uuid.Parse(cookie)
		if err != nil {
			c.JSON(http.StatusUnauthorized, apierr.CookieNotExists)
			return
		}

		session, err := sessionRepo.GetById(cookieIdUUID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, apierr.InternalServerError)
			logger.Error(err.Error())
			return
		}

		if session == nil {
			c.JSON(http.StatusUnauthorized, apierr.SessionNotFound)
			return
		}

		u, err := userRepo.GetById(session.UserId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, apierr.InternalServerError)
			logger.Error(err.Error())
			return
		}

		if u == nil {
			c.JSON(http.StatusUnauthorized, apierr.UserNotFound)
			return
		}

		c.JSON(http.StatusCreated, mapUserToUserOut(u))
	})
}

func mapUserToUserOut(u *user.User) *user.UserOut {
	return &user.UserOut{
		Id:        u.Id,
		CreatedAt: u.CreatedAt,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Username,
		Role:      u.Role,
	}
}
