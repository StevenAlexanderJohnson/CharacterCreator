package services

import (
	"dndcc/internal/models"
	"dndcc/internal/repositories"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/StevenAlexanderJohnson/grove"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo          *repositories.AuthRepository
	authenticator *grove.Authenticator[*models.Claims]
}

func NewAuthService(repo *repositories.AuthRepository, authenticator *grove.Authenticator[*models.Claims]) *AuthService {
	return &AuthService{repo: repo, authenticator: authenticator}
}

func validatePasswordRequirements(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	const specialChars = "!@#$%^&*()_+-=[]{}|;':\",.<>/?`~"

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	// 2. Check all requirements
	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func (s *AuthService) Create(data *models.Auth) (*models.Auth, error) {
	if strings.TrimSpace(data.Username) == "" {
		return nil, fmt.Errorf("username is required")
	}
	if err := validatePasswordRequirements(data.Password); err != nil {
		return nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while generating password hash: %v", err)
	}
	data.HashedPassword = string(hashedPassword)
	return s.repo.Create(data)
}

func (s *AuthService) Get(user *models.Auth) (*models.Auth, string, error) {
	data, err := s.repo.Get(user.Username)
	if err != nil {
		return nil, "", fmt.Errorf("an error occurred while getting the user in auth service: %v", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(data.HashedPassword), []byte(user.Password))
	if err != nil {
		return nil, "", fmt.Errorf("user provided invalid password for user %d", data.ID)
	}

	token, err := s.authenticator.GenerateToken(&models.Claims{
		Id:       strconv.Itoa(data.ID),
		Username: data.Username,
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    s.authenticator.Issuer,
			Subject:   data.Username,
			Audience:  s.authenticator.Audience,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.authenticator.Lifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to create token for user %d", data.ID)
	}

	return data, token, nil
}
