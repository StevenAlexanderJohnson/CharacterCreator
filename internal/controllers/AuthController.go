package controllers

import (
	"dndcc/internal/models"
	"dndcc/internal/models/page"
	"dndcc/internal/services"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/StevenAlexanderJohnson/grove"
)

type AuthController struct {
	authService    *services.AuthService
	sessionService *services.SessionService
	logger         grove.ILogger
	pageTemplates  map[string]*template.Template
}

func NewAuthController(service *services.AuthService, sessionService *services.SessionService, logger grove.ILogger) *AuthController {
	pageTemplates := make(map[string]*template.Template)

	pageTemplates["login"] = template.Must(template.ParseFiles(
		"internal/templates/layouts/layout.html.tmpl",
		"internal/templates/pages/auth.html.tmpl",
	))

	pageTemplates["register"] = template.Must(template.ParseFiles(
		"internal/templates/layouts/layout.html.tmpl",
		"internal/templates/pages/register.html.tmpl",
	))
	return &AuthController{
		authService:    service,
		sessionService: sessionService,
		logger:         logger,
		pageTemplates:  pageTemplates,
	}
}

func (c *AuthController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /auth/login", c.LoginPage)
	mux.HandleFunc("GET /auth/register", c.RegisterPage)
	mux.HandleFunc("GET /auth/logout", c.Logout)
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
	pageData := page.NewPageData(false, nil, page.LoginData{

		SessionId: session,
		Error:     "",
	})
	if err := c.pageTemplates["login"].Execute(w, pageData); err != nil {
		c.logger.Error("an error occurred while rendering login page")
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, "")
		return
	}
}

func (c *AuthController) RegisterPage(w http.ResponseWriter, r *http.Request) {
	pageData := page.NewPageData(false, nil, &page.RegisterPageData{})
	if err := c.pageTemplates["register"].Execute(w, pageData); err != nil {
		c.logger.Error("an error occurred while rendering register page")
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, "")
		return
	}
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	clearAuthCookies(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.logger.Warning("received an invalid form request to the login page: %s", r.RemoteAddr)
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "Failed to parse form")
		return
	}
	var data models.Auth
	data.SessionToken = r.FormValue("SessionId")
	data.Username = r.FormValue("Username")
	data.Password = r.FormValue("Password")
	user, authToken, duration, err := c.authService.Get(&data)
	if err != nil {
		c.logger.Error(err.Error())
		pageData := page.NewPageData(false, nil, page.LoginData{

			SessionId: data.SessionToken,
			Error:     "We could not log you in, please check your credentials and try again.",
		})
		if err := c.pageTemplates["login"].Execute(w, &pageData); err != nil {
			c.logger.Error(err.Error())
			grove.WriteErrorToResponse(w, http.StatusInternalServerError, "")
		}
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

	SetAuthCookie(w, authToken, duration)
	SetSessionCookie(w, data.SessionToken, time.Duration(24*time.Hour))

	http.Redirect(w, r, "/", http.StatusSeeOther)
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
