package services

import (
	"dndcc/internal/models"
	"dndcc/internal/repositories"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/StevenAlexanderJohnson/grove"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	repo          *repositories.AuthRepository
	authenticator *grove.Authenticator[*models.Claims]
}

func NewAuthService(repo *repositories.AuthRepository, authenticator *grove.Authenticator[*models.Claims]) *AuthService {
	return &AuthService{repo: repo, authenticator: authenticator}
}

func (s *AuthService) generateToken(user *models.Auth) (string, error) {
	return s.authenticator.GenerateToken(&models.Claims{
		UserId:   user.ID,
		Username: user.Username,
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    s.authenticator.Issuer,
			Subject:   user.Username,
			Audience:  s.authenticator.Audience,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.authenticator.Lifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})
}

func (s *AuthService) Get(user *models.Auth) (*models.Auth, string, time.Duration, error) {
	data, err := s.repo.Get(user.Username)
	if err != nil {
		return nil, "", 0, fmt.Errorf("an error occurred while getting the user in auth service: %v", err)
	}

	token, err := s.generateToken(data)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to create token for user %d: %v", data.ID, err)
	}

	return data, token, s.authenticator.Lifetime, nil
}

func (s *AuthService) GetTokenById(userId int) (string, time.Duration, error) {
	data, err := s.repo.GetId(userId)
	if err != nil {
		return "", 0, fmt.Errorf("an error occurred while getting user's data by id")
	}

	token, err := s.generateToken(data)
	if err != nil {
		return "", 0, fmt.Errorf("failed tdo create token when getting user by id: %v", err)
	}

	return token, s.authenticator.Lifetime, nil
}

func (s *AuthService) ValidateOAuth2(token, token_id string) (*models.Auth, string, time.Duration, error) {
	keys := []models.OAuth2PublicKey{}
	resp, err := http.Get("http://localhost:8081/.well-known/jwks.json")
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to get well known jwks.json: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to read the well known jwks.json response body: %v", err)
	}

	if err := json.Unmarshal(body, &keys); err != nil {
		return nil, "", 0, fmt.Errorf("failed to unmarshal well known jwks.json: %v", err)
	}

	var key *models.OAuth2PublicKey
	for _, k := range keys {
		if k.ID == token_id {
			key = &k
			break
		}
	}
	if key == nil {
		return nil, "", 0, fmt.Errorf("failed to find valid key in well known jwks.json response")
	}

	oauth2Claims := &models.OAuth2Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, oauth2Claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return &key.PublicKey, nil
	})
	if err != nil {
		return nil, "", 0, fmt.Errorf("unable to parse OAuth2 token: %v", err)
	}
	if !parsedToken.Valid {
		return nil, "", 0, fmt.Errorf("the OAuth2 claims validation failed")
	}

	user, err := s.repo.Get(oauth2Claims.Username)
	if err != nil {
		if !errors.Is(err, repositories.ErrUserNotFound) {
			return nil, "", 0, fmt.Errorf("an error occurred while finding OAuth2 user in the database: %v", err)
		}
		user, err = s.repo.Create(oauth2Claims)
		if err != nil {
			return nil, "", 0, fmt.Errorf("failed to adding new OAuth2 user to database: %v", err)
		}
	}

	authToken, err := s.generateToken(user)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to generate token for OAuth2 user: %v", err)
	}

	return user, authToken, s.authenticator.Lifetime, nil
}
