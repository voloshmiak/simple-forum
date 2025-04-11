package render

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type Renderer struct {
	templateCache map[string]*template.Template
}

func NewRenderer() (*Renderer, error) {
	r := &Renderer{
		templateCache: make(map[string]*template.Template),
	}
	err := r.InitTemplates()
	return r, err
}

func (r *Renderer) InitTemplates() error {
	tc, err := r.createTemplateCache()
	r.templateCache = tc
	return err
}

func (r *Renderer) RenderTemplate(rw http.ResponseWriter, tmpl string, data any) error {
	// if in debug mode
	if os.Getenv("DEBUG_MODE") == "true" {
		err := r.InitTemplates()
		if err != nil {
			return err
		}
	}

	// get requested template
	rt, ok := r.templateCache[tmpl+".gohtml"]
	if !ok {
		http.Error(rw, tmpl+".gohtml not found", http.StatusNotFound)
		return errors.New(tmpl + ".gohtml not found")
	}

	// render template
	err := rt.Execute(rw, data)
	if err != nil {
		http.Error(rw, tmpl+".gohtml failed to render", http.StatusInternalServerError)
		return err
	}

	return nil
}

func (r *Renderer) createTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// getting path to templates
	templatesPath := filepath.Join("internal", "web", "templates")

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
