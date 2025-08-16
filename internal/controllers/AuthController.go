package controllers

import (
	"dndcc/internal/models"
	"dndcc/internal/services"
	"fmt"
	"net/http"

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
	mux.HandleFunc("POST /auth/login", c.Login)
	mux.HandleFunc("POST /auth/register", c.Register)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	data, err := grove.ParseJsonBodyFromRequest[*models.Auth](r)
	if err != nil {
		c.logger.Error(err.Error())
		grove.WriteErrorToResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	user, authToken, err := c.authService.Get(data)
	if err != nil {
		c.logger.Error(err.Error())
		grove.WriteErrorToResponse(w, http.StatusNoContent, "")
		return
	}
	if data.SessionToken != "" {
		session, err := c.sessionService.Get(data.SessionToken)
		if err != nil {
			c.logger.Warning(fmt.Errorf("a log attempt for %s was attempted with invalid session: %v", data.SessionToken, err))
			grove.WriteErrorToResponse(w, http.StatusNoContent, "")
			return
		}
		if session.UserId != user.ID {
			c.logger.Warning(fmt.Errorf("a log attempt for %s with session token was made from an invalid user %s: %v", user.Username, data.SessionToken, err))
			grove.WriteErrorToResponse(w, http.StatusNoContent, "")
			return
		}
	} else {
		session, err := c.sessionService.Create(models.CreateNewSession(user.ID, r.UserAgent(), r.RemoteAddr))
		if err != nil {
			c.logger.Warning(fmt.Errorf("an error occurred while creating a new session: %v", err))
			grove.WriteErrorToResponse(w, http.StatusInternalServerError, "")
			return
		}
		data.SessionToken = session.Token
	}

	if err := grove.WriteJsonBodyToResponse(
		w,
		struct {
			AuthToken    string `json:"auth_token"`
			SessionToken string `json:"session_token"`
		}{
			AuthToken:    authToken,
			SessionToken: data.SessionToken,
		},
	); err != nil {
		c.logger.Error(err.Error())
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	data, err := grove.ParseJsonBodyFromRequest[*models.Auth](r)
	if err != nil {
		grove.WriteErrorToResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	newData, err := c.authService.Create(data)
	if err != nil {
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := grove.WriteJsonBodyToResponse(w, newData); err != nil {
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
