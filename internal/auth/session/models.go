package session

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	UserId    int64     `json:"userId"`
}

func (session *Session) IsValid() bool {
	return !time.Now().After(session.ExpiresAt)
}
