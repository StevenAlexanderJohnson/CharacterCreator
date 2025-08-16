package models

import (
	"crypto/rand"
	"time"
)

type Session struct {
	ID             int
	CreatedAt      time.Time
	ExpiresAt      time.Time
	IpAddress      string
	LastActivityAt time.Time
	Token          string
	UserAgent      string
	UserId         int
}

func CreateNewSession(userId int, userAgent, ip string) *Session {
	return &Session{
		CreatedAt:      time.Now(),
		ExpiresAt:      time.Now(),
		IpAddress:      ip,
		LastActivityAt: time.Now(),
		Token:          rand.Text(),
		UserAgent:      userAgent,
		UserId:         userId,
	}
}
