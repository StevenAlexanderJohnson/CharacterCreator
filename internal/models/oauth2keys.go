package models

import "crypto/rsa"

type OAuth2PublicKey struct {
	ID        string `json:"id"`
	PublicKey rsa.PublicKey
}
