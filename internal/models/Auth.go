package models

import (
	"time"
)

type Auth struct {
	ID           int       `json:"-"`
	Username     string    `json:"username"`
	SessionToken string    `json:"session_token,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func AuthFromOAuth2(jwt *OAuth2Claims) *Auth {
	return &Auth{
		ID:       jwt.ID,
		Username: jwt.Username,
	}
}
