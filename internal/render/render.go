package render

import (
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"path/filepath"
)

var templatesPath string = "../../internal/web/templates/"

type Renderer struct {
	templateCache map[string]*template.Template
	logger        *zap.SugaredLogger
	Debug         bool
}

func NewRenderer(logger *zap.SugaredLogger) *Renderer {
	r := &Renderer{
		templateCache: make(map[string]*template.Template),
		logger:        logger,
	}
	r.InitTemplates()
	return r
}

func (r *Renderer) InitTemplates() {
	tc, err := r.createTemplateCache()
	r.templateCache = tc
	if err != nil {
		r.logger.Fatal("init templates failed", "err", err)
	}
}

func (r *Renderer) RenderTemplate(rw http.ResponseWriter, tmpl string, data any) {
	// if in debug mode
	if r.Debug {
		r.InitTemplates()
	}

	// get requested template
	rt, ok := r.templateCache[tmpl+".gohtml"]
	if !ok {
		http.Error(rw, tmpl+".gohtml not found", http.StatusNotFound)
		r.logger.Error(tmpl + ".gohtml not found")
		return
	}

	// render template
	err := rt.Execute(rw, data)
	if err != nil {
		http.Error(rw, tmpl+".gohtml failed to render", http.StatusInternalServerError)
		r.logger.Error(tmpl+".gohtml failed to render", "err", err)
		return
	}
}

func (r *Renderer) createTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(templatesPath + "*.page.gohtml")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		tmpl, err := template.New(name).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		matches, err := filepath.Glob(templatesPath + "*.layout.gohtml")

		if err != nil {
			return nil, err
		}

		if len(matches) > 0 {
			tmpl, err = tmpl.ParseGlob(templatesPath + "*.layout.gohtml")
			if err != nil {
				return nil, err
			}
		}

		myCache[name] = tmpl
	}

	return myCache, nil
}
