package template

import (
	"errors"
	"forum-project/internal/auth"
	"forum-project/internal/config"
	"forum-project/internal/model"
	"html/template"
	"net/http"
	"path/filepath"
)

type Templates struct {
	Cache  map[string]*template.Template
	Config *config.Config
}

func NewTemplates(config *config.Config) *Templates {
	cache := parseTemplates(config.Path.ToTemplates())
	return &Templates{
		Cache:  cache,
		Config: config,
	}
}

func parseTemplates(basePath string) map[string]*template.Template {
	templates := map[string]*template.Template{}

	// getting path to templates
	templatesPath := basePath

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

func (m *Templates) AddDefaultData(td *model.Page, r *http.Request) *model.Page {
	claims, err := auth.GetClaimsFromRequest(r, m.Config.JWT.Secret)
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

func (m *Templates) Render(rw http.ResponseWriter, r *http.Request, tmpl string, td *model.Page) error {
	// cache if in development mode
	if m.Config.Env == "development" {
		templates := parseTemplates(m.Config.Path.ToTemplates())
		m.Cache = templates
	}

	// get requested template
	rt, ok := m.Cache[tmpl+".gohtml"]
	if !ok {
		return errors.New(tmpl + ".gohtml not found")
	}

	td = m.AddDefaultData(td, r)

	// rendering template
	return rt.Execute(rw, td)
}
