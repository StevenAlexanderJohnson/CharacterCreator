package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Id string `json:"id"`
	*jwt.RegisteredClaims
}
