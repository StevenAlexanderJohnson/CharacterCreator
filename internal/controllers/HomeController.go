package controllers

import (
	"dndcc/internal/models"
	"dndcc/internal/models/page"
	"html/template"
	"net/http"

	"github.com/StevenAlexanderJohnson/grove"
)

type HomeController struct {
	logger        grove.ILogger
	authenticator *grove.Authenticator[*models.Claims]
	pageTemplates map[string]*template.Template
}

func NewHomeController(logger grove.ILogger, authenticator *grove.Authenticator[*models.Claims]) *HomeController {
	pageTemplates := make(map[string]*template.Template)
	pageTemplates["index"] = template.Must(template.ParseFiles(
		"internal/templates/layouts/layout.html.tmpl",
		"internal/templates/pages/home.html.tmpl",
	))

	return &HomeController{
		logger,
		authenticator,
		pageTemplates,
	}
}

func (h *HomeController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", h.Index)
	mux.HandleFunc("GET /health", h.HealthCheck)
}

func (h *HomeController) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	claims, ok := r.Context().Value(grove.AuthTokenKey).(*models.Claims)

	pageData := page.PageData[page.HomePageData]{
		IsAuthenticated: ok,
		User:            claims,
		Data:            page.HomePageData{Authenticated: ok},
	}

	if err := h.pageTemplates["index"].Execute(w, pageData); err != nil {
		h.logger.Error("an error occurred while rendering home page")
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, "")
		return
	}
}

func (h *HomeController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
