package template

import (
	"errors"
	"forum-project/internal/auth"
	"forum-project/internal/model"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
)

type Templates struct {
	cache  map[string]*template.Template
	env    string
	path   string
	auther *auth.JwtAuthenticator
}

func NewTemplates(env, path string, auther *auth.JwtAuthenticator) *Templates {
	cache := parseTemplates(path)
	return &Templates{
		cache:  cache,
		env:    env,
		path:   path,
		auther: auther,
	}
}

func parseTemplates(basePath string) map[string]*template.Template {
	templates := map[string]*template.Template{}

	// getting path to templates
	templatesPath := basePath

	// parsing templates
	layouts, _ := filepath.Glob(filepath.Join(templatesPath, "*.layout.gohtml"))

	pages, _ := filepath.Glob(filepath.Join(templatesPath, "*.page.gohtml"))

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

func (m *Templates) addDefaultData(td *model.Page, r *http.Request) *model.Page {
	td.CSRFToken = nosurf.Token(r)
	td.IsAuthenticated = false
	td.IsAdmin = false

	claims, err := m.auther.GetClaimsFromRequest(r)
	if err != nil {
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
	if m.env == "development" {
		templates := parseTemplates(m.path)
		m.cache = templates
	}

	// get requested template
	rt, ok := m.cache[tmpl+".gohtml"]
	if !ok {
		return errors.New(tmpl + ".gohtml not found")
	}

	td = m.addDefaultData(td, r)

	// rendering template
	return rt.Execute(rw, td)
}
