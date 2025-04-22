package template

import (
	"errors"
	"html/template"
	"net/http"
	"path/filepath"
)

type Manager struct {
	templates map[string]*template.Template
	debugMode bool
}

func NewManager(debugMode bool) (*Manager, error) {
	templates, err := parseTemplates()
	if err != nil {
		return nil, err
	}
	r := &Manager{
		templates: templates,
		debugMode: debugMode,
	}
	return r, nil
}

func (m *Manager) Render(rw http.ResponseWriter, tmpl string, data any) error {
	// if in debug mode
	if m.debugMode {
		templates, err := parseTemplates()
		if err != nil {
			return nil
		}
		m.templates = templates
	}

	// get requested template
	rt, ok := m.templates[tmpl+".gohtml"]
	if !ok {
		http.Error(rw, tmpl+".gohtml not found", http.StatusNotFound)
		return errors.New(tmpl + ".gohtml not found")
	}

	// rendering template
	return rt.Execute(rw, data)
}

func parseTemplates() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// getting path to templates
	templatesPath := filepath.Join("web", "templates")

	pages, err := filepath.Glob(templatesPath + "\\*.page.gohtml")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		tmpl, err := template.New(name).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		matches, err := filepath.Glob(templatesPath + "\\*.layout.gohtml")

		if err != nil {
			return nil, err
		}

		if len(matches) > 0 {
			tmpl, err = tmpl.ParseGlob(templatesPath + "\\*.layout.gohtml")
			if err != nil {
				return nil, err
			}
		}

		myCache[name] = tmpl
	}

	return myCache, nil
}
