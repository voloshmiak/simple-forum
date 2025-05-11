package template

import (
	"errors"
	"forum-project/internal/auth"
	"forum-project/internal/models"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type Templates struct {
	cache map[string]*template.Template
}

func NewTemplates() (*Templates, error) {
	templates, err := parseTemplates()
	if err != nil {
		return nil, err
	}
	return &Templates{cache: templates}, nil
}

func addDefaultData(td *models.Page, r *http.Request) *models.Page {
	claims, err := auth.GetClaimsFromRequest(r)
	if err != nil {
		td.IsAuthenticated = false
		td.IsAdmin = false
		return td
	}

	td.IsAuthenticated = true
	td.IsAdmin = false

	user := claims["user"].(map[string]interface{})

	role := user["role"].(string)

	if role == "admin" {
		td.IsAdmin = true
	}

	userName := user["username"].(string)

	stringMap := make(map[string]string)
	stringMap["username"] = userName

	td.StringMap = stringMap

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

func (m *Templates) Render(rw http.ResponseWriter, r *http.Request, tmpl string, td *models.Page) error {
	// if in development mode
	isDevelopment := os.Getenv("APP_ENV") == "development"
	if isDevelopment {
		templates, err := parseTemplates()
		if err != nil {
			return err
		}
		m.cache = templates
	}

	// get requested template
	rt, ok := m.cache[tmpl+".gohtml"]
	if !ok {
		return errors.New(tmpl + ".gohtml not found")
	}

	td = addDefaultData(td, r)

	// rendering template
	return rt.Execute(rw, td)
}
