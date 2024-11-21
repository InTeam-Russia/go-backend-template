package auth

import (
	"github.com/InTeam-Russia/go-backend-template/internal/auth/api"
	"github.com/InTeam-Russia/go-backend-template/internal/auth/session"
	"github.com/InTeam-Russia/go-backend-template/internal/auth/user"
)

type CookieConfig = api.CookieConfig
type UserRepository = user.UserRepository
type CreateUser = user.CreateUser
type SessionRepository = session.SessionRepository

var SetupRoutes = api.SetupRoutes
var DefaultCookieConfig = api.DefaultCookieConfig
var NewPgUserRepository = user.NewPgUserRepository
var NewRedisSessionRepository = session.NewRedisSessionRepository
