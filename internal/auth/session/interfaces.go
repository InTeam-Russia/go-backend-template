package session

import "github.com/google/uuid"

type Seconds = int

type SessionRepository interface {
	Create(userId int64, sessionLifetime Seconds) (*Session, error)
	GetById(id uuid.UUID) (*Session, error)
	DeleteById(id uuid.UUID) error
}
