package controllers

import (
	"dndcc/internal/models"
	"dndcc/internal/services"
	"net/http"
	"time"

	"github.com/StevenAlexanderJohnson/grove"
)

type AuthController struct {
	authService    *services.AuthService
	sessionService *services.SessionService
	logger         grove.ILogger
}

func NewAuthController(service *services.AuthService, sessionService *services.SessionService, logger grove.ILogger) *AuthController {
	return &AuthController{
		authService:    service,
		sessionService: sessionService,
		logger:         logger,
	}
}

func (c *AuthController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /auth/login", c.Login)
	mux.HandleFunc("GET /auth/register", c.Register)
	mux.HandleFunc("GET /auth/logout", c.Logout)
	mux.HandleFunc("GET /auth/validate", c.ValidateOAuth2)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:8081/login?service_name=localhost:8080&redirect_uri=/auth/validate", http.StatusSeeOther)
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:8081/register?service_name=localhost:8080&redirect_uri=/auth/validate", http.StatusSeeOther)
}

func (c *AuthController) ValidateOAuth2(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	token_id := r.URL.Query().Get("token_id")
	if token == "" || token_id == "" {
		c.logger.Error("a request was made to validate OAuth2 with invalid empty token or token_id")
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, "")
		return
	}

	user, token, duration, err := c.authService.ValidateOAuth2(token, token_id)
	if err != nil {
		c.logger.Errorf("an error occurred while validating the OAuth2 response: %v", err)
		grove.WriteErrorToResponse(w, http.StatusNoContent, "")
		return
	}
	session, err := c.sessionService.Create(models.CreateNewSession(user.ID, r.UserAgent(), r.RemoteAddr))
	if err != nil {
		c.logger.Errorf("an error occurred while creating a new session: %v", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, "")
		return
	}

	SetAuthCookie(w, token, duration)
	SetSessionCookie(w, session.Token, time.Duration(32*time.Hour))

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	clearAuthCookies(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
