package models

import "time"

type Auth struct {
	ID              int       `json:"-"`
	ConfirmPassword string    `json:"confirm_password,omitempty"`
	Password        string    `json:"password,omitempty"`
	Username        string    `json:"username"`
	HashedPassword  string    `json:"-"`
	SessionToken    string    `json:"session_token,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
