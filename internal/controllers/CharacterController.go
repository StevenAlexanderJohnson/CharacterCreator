package controllers

import (
	"database/sql"
	"dndcc/internal/character"
	"dndcc/internal/models"
	"dndcc/internal/models/page"
	"dndcc/internal/services"
	"fmt"
	"html/template"
	"net/http"
	"slices"
	"strconv"

	"github.com/StevenAlexanderJohnson/grove"
)

type CharacterController struct {
	logger        grove.ILogger
	service       *services.CharacterService
	pageTemplates map[string]*template.Template
}

func NewCharacterController(logger grove.ILogger, service *services.CharacterService) *CharacterController {
	pageTemplates := make(map[string]*template.Template)
	funcMap := template.FuncMap{
		"statCard": func(name string, score int, modifier int) map[string]interface{} {
			return map[string]interface{}{
				"Name":     name,
				"Score":    score,
				"Modifier": modifier,
			}
		},
		"skill": func(name string, char *character.Character, base character.StatName) map[string]interface{} {
			skillName := character.SkillName(name)
			hasProficiency := func() bool {
				for _, prof := range char.Background.GetProficiencies() {
					if prof == skillName {
						return true
					}
				}
				return false
			}()
			return map[string]interface{}{
				"Name":           name,
				"Bonus":          char.GetSkill(skillName),
				"HasProficiency": hasProficiency,
				"StatName":       base,
			}
		},
		"savingThrow": func(stat character.StatName, char *character.Character) map[string]interface{} {
			hasProficiency := func() bool {
				for _, prof := range character.ClassBarbarian.GetSavingThrowsProficiencies() {
					if prof == stat {
						return true
					}
				}
				return false
			}()
			bonus := char.GetSavingThrow(stat)
			return map[string]interface{}{
				"Name":           stat,
				"HasProficiency": hasProficiency,
				"Bonus":          bonus,
			}
		},
	}

	pageTemplates["list"] = template.Must(template.ParseFiles(
		"internal/templates/layouts/layout.html.tmpl",
		"internal/templates/pages/characterList.html.tmpl",
	))

	pageTemplates["new"] = template.Must(template.New("edit").Funcs(template.FuncMap{
		"isCustomBackground": func(list []character.BackgroundName, item string) bool {
			return !slices.Contains(list, character.BackgroundName(item))
		},
		"isCustomRace": func(list []character.RaceName, item string) bool {
			return !slices.Contains(list, character.RaceName(item))
		},
		"isCustomSubrace": func(list []character.SubraceName, item sql.NullString) bool {
			if !item.Valid {
				return false
			}
			return !slices.Contains(list, character.SubraceName(item.String))
		},
	}).ParseFiles(
		"internal/templates/layouts/layout.html.tmpl",
		"internal/templates/pages/characterEdit.html.tmpl",
	))

	pageTemplates["character"] = template.Must(template.New("character").Funcs(funcMap).ParseFiles(
		"internal/templates/layouts/layout.html.tmpl",
		"internal/templates/partials/statCard.html.tmpl",
		"internal/templates/partials/skill.html.tmpl",
		"internal/templates/partials/savingThrow.html.tmpl",
		"internal/templates/pages/character.html.tmpl",
	))

	return &CharacterController{
		logger:        logger,
		service:       service,
		pageTemplates: pageTemplates,
	}
}

func (c *CharacterController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /character", c.Create)
	mux.HandleFunc("GET /character", c.GetAll)
	mux.HandleFunc("GET /character/new", c.NewCharacter)
	mux.HandleFunc("GET /character/{id}", c.GetByID)
	mux.HandleFunc("GET /character/{id}/edit", c.EditCharacter)
	mux.HandleFunc("PUT /character/{id}", c.Update)
	mux.HandleFunc("DELETE /character/{id}", c.Delete)
}

func (c *CharacterController) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(grove.AuthTokenKey).(*models.Claims)
	if !ok {
		c.logger.Warning("unauthenticated user reached /character/new endpoint")
		grove.WriteErrorToResponse(w, http.StatusUnauthorized, "")
		return
	}

	if err := r.ParseForm(); err != nil {
		c.logger.Warning("an invalid request reached /character/new endpoint")
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "failed to parse form")
		return
	}

	data, err := models.CharacterFromForm(r)
	if err != nil {
		c.logger.Warning("parsing request form to character in /character/new endpoint failed: %v", err)
		if err := c.pageTemplates["new"].ExecuteTemplate(w, "layout.html.tmpl", page.NewCharacterEditPageData("/character", "post", err.Error(), data)); err != nil {
			c.logger.Error("an error occurred while rendering the edit page after failed create", err)
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	data.OwnerId = claims.UserId
	_, err = c.service.Create(data)
	if err != nil {
		c.logger.Error("an error occurred while updating character", err)
		pageData := page.NewCharacterEditPageData("post", "/character", err.Error(), data)
		if err := c.pageTemplates["new"].ExecuteTemplate(w, "content", pageData); err != nil {
			c.logger.Error("an error occurred while rendering the edit page after failed update", err)
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/character/%d", data.ID))
}

func (c *CharacterController) GetAll(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(grove.AuthTokenKey).(*models.Claims)
	if !ok {
		grove.WriteErrorToResponse(w, http.StatusUnauthorized, "")
		return
	}

	items, err := c.service.List(claims.UserId)
	if err != nil {
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	pageData := page.NewPageData(ok, claims, items)

	if err := c.pageTemplates["list"].Execute(w, pageData); err != nil {
		c.logger.Error("failed to render template list within the character controller", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (c *CharacterController) NewCharacter(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(grove.AuthTokenKey).(*models.Claims)
	pageData := page.NewPageData(ok, claims, page.NewCharacterEditPageData(
		"post",
		"/character",
		"",
		&models.Character{Level: 1},
	))
	if err := c.pageTemplates["new"].ExecuteTemplate(w, "layout.html.tmpl", pageData); err != nil {
		c.logger.Error("failed to render template new within the character controller", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (c *CharacterController) GetByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(grove.AuthTokenKey).(*models.Claims)
	if !ok {
		grove.WriteErrorToResponse(w, http.StatusUnauthorized, "")
		return
	}

	idString := r.PathValue("id")
	if idString == "" {
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "ID is required")
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	item, err := c.service.Get(id, claims.UserId)
	if err != nil {
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if item == nil {
		grove.WriteErrorToResponse(w, http.StatusNotFound, "Item not found")
		return
	}

	pageData := page.NewPageData(ok, claims, page.NewCharacterViewPageData(item.ID, item.ToCharacterSheet()))
	if err := c.pageTemplates["character"].ExecuteTemplate(w, "layout.html.tmpl", pageData); err != nil {
		c.logger.Error("an error occurred while rendering character page", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, "")
		return
	}
}

func (c *CharacterController) EditCharacter(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(grove.AuthTokenKey).(*models.Claims)
	if !ok {
		grove.WriteErrorToResponse(w, http.StatusUnauthorized, "")
		return
	}

	idString := r.PathValue("id")
	if idString == "" {
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "ID is required")
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}
	item, err := c.service.Get(id, claims.UserId)
	if err != nil {
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if item == nil {
		grove.WriteErrorToResponse(w, http.StatusNotFound, "Item not found")
		return
	}
	pageData := page.NewPageData(ok, claims, page.NewCharacterEditPageData(
		"put",
		fmt.Sprintf("/character/%d", item.ID),
		"",
		item,
	))
	if err := c.pageTemplates["new"].ExecuteTemplate(w, "layout.html.tmpl", pageData); err != nil {
		c.logger.Error("failed to render template edit within the character controller", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (c *CharacterController) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(grove.AuthTokenKey).(*models.Claims)
	if !ok {
		c.logger.Warning("unauthenticated user reached /character/new endpoint")
		grove.WriteErrorToResponse(w, http.StatusUnauthorized, "")
		return
	}

	if err := r.ParseForm(); err != nil {
		c.logger.Warning("an invalid request reached /character/new endpoint")
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "failed to parse form")
		return
	}

	data, err := models.CharacterFromForm(r)
	if err != nil {
		c.logger.Warning("parsing request form to character in /character/new endpoint failed: %v", err)
		if err := c.pageTemplates["new"].ExecuteTemplate(w, "layout.html.tmpl", page.NewCharacterEditPageData("/character", "post", err.Error(), data)); err != nil {
			c.logger.Error("an error occurred while rendering the edit page after failed create", err)
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	idString := r.PathValue("id")
	if idString == "" {
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "ID is required")
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	updatedData, err := c.service.Update(data, id, claims.UserId)
	if err != nil {
		c.logger.Error("an error occurred while updating character", err)
		pageData := page.NewCharacterEditPageData("put", fmt.Sprintf("/character/%d", id), err.Error(), data)
		if err := c.pageTemplates["new"].ExecuteTemplate(w, "content", pageData); err != nil {
			c.logger.Error("an error occurred while rendering the edit page after failed update", err)
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/character/%d", updatedData.ID))
}

func (c *CharacterController) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(grove.AuthTokenKey).(*models.Claims)
	if !ok {
		grove.WriteErrorToResponse(w, http.StatusUnauthorized, "")
		return
	}

	idString := r.PathValue("id")
	if idString == "" {
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "ID is required")
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if err := c.service.Delete(id, claims.UserId); err != nil {
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
