package controllers

import (
	"dndcc/internal/models"
	"dndcc/internal/services"
	"net/http"

	"github.com/StevenAlexanderJohnson/grove"
)

type AuthController struct {
	service *services.AuthService
	logger  grove.ILogger
}

func NewAuthController(service *services.AuthService, logger grove.ILogger) *AuthController {
	return &AuthController{
		service: service,
		logger:  logger,
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
	authToken, err := c.service.Get(data)
	if err != nil {
		c.logger.Error(err.Error())
		grove.WriteErrorToResponse(w, http.StatusNoContent, "")
		return
	}

	if err := grove.WriteJsonBodyToResponse(
		w,
		struct {
			AuthToken string `json:"auth_token"`
		}{
			AuthToken: authToken,
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
	newData, err := c.service.Create(data)
	if err != nil {
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := grove.WriteJsonBodyToResponse(w, newData); err != nil {
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
