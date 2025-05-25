package template

import (
	"errors"
	"forum-project/internal/auth"
	"forum-project/internal/env"
	"forum-project/internal/model"
	"html/template"
	"net/http"
	"path/filepath"
)

type Renderer interface {
	Render(rw http.ResponseWriter, r *http.Request, tmpl string, td *model.Page) error
}

type Templates struct {
	cache map[string]*template.Template
}

func NewTemplates() *Templates {
	return &Templates{
		cache: parseTemplates(),
	}
}

func addDefaultData(td *model.Page, r *http.Request) *model.Page {
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

func parseTemplates() map[string]*template.Template {
	templates := map[string]*template.Template{}

	// getting path to templates
	templatesPath := env.GetTemplatePath()

	// parsing templates
	layouts, _ := filepath.Glob(templatesPath + "\\*.layout.gohtml")

	pages, _ := filepath.Glob(templatesPath + "\\*.page.gohtml")

	for _, page := range pages {
		name := filepath.Base(page)

		filenames := make([]string, 0, len(layouts)+1)

		filenames = append(filenames, page)
		filenames = append(filenames, layouts...)

		tmpl, _ := template.New(name).ParseFiles(filenames...)

		templates[name] = tmpl
	}

	return templates
}

func (m *Templates) Render(rw http.ResponseWriter, r *http.Request, tmpl string, td *model.Page) error {
	// if in development mode
	isDevelopment := env.GetEnv("APP_ENV", "development") == "development"
	if isDevelopment {
		templates := parseTemplates()
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
