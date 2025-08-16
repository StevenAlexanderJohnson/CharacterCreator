package controllers

import (
	"dndcc/internal/models"
	"dndcc/internal/models/page"
	"dndcc/internal/services"
	"fmt"
	"html/template"
	"net/http"

	"github.com/StevenAlexanderJohnson/grove"
)

var authTemplates = make(map[string]*template.Template)

func init() {
	authTemplates["login"] = template.Must(template.ParseFiles(
		"internal/templates/layouts/layout.html.tmpl",
		"internal/templates/pages/auth.html.tmpl",
	))
}

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
	mux.HandleFunc("GET /auth/login", c.LoginPage)
	mux.HandleFunc("POST /auth/login", c.Login)
	mux.HandleFunc("POST /auth/register", c.Register)
}

func (c *AuthController) LoginPage(w http.ResponseWriter, r *http.Request) {
	session, err := getSessionCookie(r)
	if err != nil {
		c.logger.Warning("an error occurred while trying to get session id for login page", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, "")
		return
	}
	if err := authTemplates["login"].Execute(w, page.LoginData{
		SessionId: session,
		Error:     "",
	}); err != nil {
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, "")
		return
	}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.logger.Warning("received an invalid form request to the login page: %s", r.RemoteAddr)
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "Failed to parse form")
		return
	}
	var data models.Auth
	data.Username = r.FormValue("username")
	data.Password = r.FormValue("password")
	user, authToken, err := c.authService.Get(&data)
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
