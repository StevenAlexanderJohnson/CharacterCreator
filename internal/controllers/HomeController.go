package controllers

import (
	"dndcc/internal/models"
	"dndcc/internal/models/page"
	"html/template"
	"net/http"

	"github.com/StevenAlexanderJohnson/grove"
)

var homePages = make(map[string]*template.Template)

func init() {
	homePages["index"] = template.Must(template.ParseFiles(
		"internal/templates/layouts/layout.html.tmpl",
		"internal/templates/pages/home.html.tmpl",
	))
}

type HomeController struct {
	logger        grove.ILogger
	authenticator *grove.Authenticator[*models.Claims]
}

func NewHomeController(logger grove.ILogger, authenticator *grove.Authenticator[*models.Claims]) *HomeController {
	return &HomeController{
		logger,
		authenticator,
	}
}

func (h *HomeController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.Index)
}

func (h *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	token, err := getAuthCookie(r)
	authenticated := false
	if err == nil {
		authenticated = true
	}

	var claims *models.Claims = nil
	if authenticated {
		var verifyErr error
		claims, verifyErr = h.authenticator.VerifyToken(token, &models.Claims{})
		if verifyErr != nil {
			if authenticated {
				h.logger.Warning("somehow showing authenticated on home but without claims")
			}
			authenticated = false
		}
	}

	pageData := page.PageData[page.HomePageData]{
		IsAuthenticated: authenticated,
		User:            claims,
		Data:            page.HomePageData{Authenticated: authenticated},
	}

	if err := homePages["index"].Execute(w, pageData); err != nil {
		h.logger.Error("an error occurred while rendering home page")
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, "")
		return
	}
}
