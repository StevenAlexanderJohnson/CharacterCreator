package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	*jwt.RegisteredClaims
}

type OAuth2Claims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	*jwt.RegisteredClaims
}

func ClaimsFromOAuth2(oauth2Claims OAuth2Claims) *Claims {
	return &Claims{
		UserId:   oauth2Claims.ID,
		Username: oauth2Claims.Username,
	}
}
