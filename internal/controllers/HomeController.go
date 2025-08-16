package controllers

import (
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
	logger grove.ILogger
}

func NewHomeController(logger grove.ILogger) *HomeController {
	return &HomeController{
		logger,
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
	_, err := getAuthCookie(r)

	authenticated := err == nil

	if err := homePages["index"].Execute(w, &page.HomePageData{Authenticated: authenticated}); err != nil {
		h.logger.Error("an error occurred while rendering home page")
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, "")
		return
	}
}
