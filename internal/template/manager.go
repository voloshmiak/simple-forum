package template

import (
	"errors"
	"forum-project/internal/auth"
	"forum-project/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"html/template"
	"net/http"
	"path/filepath"
)

type Manager struct {
	templates map[string]*template.Template
}

func NewManager() (*Manager, error) {
	templates, err := parseTemplates()
	if err != nil {
		return nil, err
	}
	return &Manager{templates: templates}, nil
}

func addDefaultData(td *models.ViewData, r *http.Request) *models.ViewData {
	token, err := auth.ValidateTokenFromRequest(r)
	if err != nil {
		td.IsAuthenticated = false
		td.IsAdmin = false
		return td
	}

	td.IsAuthenticated = true
	td.IsAdmin = false

	claims := token.Claims.(jwt.MapClaims)

	user := claims["user"].(map[string]interface{})
	role := user["role"].(string)

	if role == "admin" {
		td.IsAdmin = true
		return td
	}

	return td
}

func parseTemplates() (map[string]*template.Template, error) {
	templates := map[string]*template.Template{}

	// getting path to templates
	templatesPath := filepath.Join("web", "templates")

	layouts, err := filepath.Glob(templatesPath + "\\*.layout.gohtml")
	if err != nil {
		return templates, err
	}

	pages, err := filepath.Glob(templatesPath + "\\*.page.gohtml")
	if err != nil {
		return templates, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		filenames := make([]string, 0, len(layouts)+1)

		filenames = append(filenames, page)
		filenames = append(filenames, layouts...)

		tmpl, err := template.New(name).ParseFiles(filenames...)
		if err != nil {
			return templates, err
		}

		templates[name] = tmpl
	}

	return templates, nil
}

func (m *Manager) Render(rw http.ResponseWriter, r *http.Request, tmpl string, td *models.ViewData) error {
	// if in development mode
	isDevelopment := true
	if isDevelopment {
		templates, err := parseTemplates()
		if err != nil {
			return err
		}
		m.templates = templates
	}

	// get requested template
	rt, ok := m.templates[tmpl+".gohtml"]
	if !ok {
		return errors.New(tmpl + ".gohtml not found")
	}

	td = addDefaultData(td, r)

	// rendering template
	return rt.Execute(rw, td)
}
