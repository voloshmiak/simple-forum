package template

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/justinas/nosurf"
	"html/template"
	"net/http"
	"path/filepath"
	"simple-forum/internal/auth"
	"simple-forum/internal/model"
)

var (
	ErrInvalidUser     = errors.New("invalid user type")
	ErrInvalidRole     = errors.New("invalid role type")
	ErrInvalidUserName = errors.New("invalid user name type")
)

type Authenticator interface {
	GetClaimsFromRequest(r *http.Request) (jwt.MapClaims, error)
}

type Templates struct {
	basePath string
	inProd   bool
	auther   Authenticator
	cache    map[string]*template.Template
}

func NewTemplates(basePath string, inProd bool, auther *auth.JWTAuthenticator) (*Templates, error) {
	templateAbsPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	templateSlashPath := filepath.ToSlash(templateAbsPath)

	cache, err := parseTemplates(templateSlashPath)
	if err != nil {
		return nil, err
	}

	return &Templates{
		cache:    cache,
		basePath: basePath,
		inProd:   inProd,
		auther:   auther,
	}, nil
}

func parseTemplates(basePath string) (map[string]*template.Template, error) {
	templates := map[string]*template.Template{}

	// parsing templates
	layouts, err := filepath.Glob(filepath.Join(basePath, "*.layout.gohtml"))
	if err != nil {
		return templates, err
	}

	pages, err := filepath.Glob(filepath.Join(basePath, "*.page.gohtml"))
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

func (m *Templates) addDefaultData(td *model.Page, r *http.Request) (*model.Page, error) {
	claims, err := m.auther.GetClaimsFromRequest(r)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return td, nil
		}
		return td, err
	}

	td.IsAuthenticated = true

	user, ok := claims["user"].(map[string]interface{})
	if !ok {
		return td, ErrInvalidUser
	}

	role, ok := user["role"].(string)
	if !ok {
		return td, ErrInvalidRole
	}

	userName, ok := user["name"].(string)
	if !ok {
		return td, ErrInvalidUserName
	}

	td.StringMap = map[string]string{
		"name": userName,
	}

	if role == "admin" {
		td.IsAdmin = true
	}

	td.CSRFToken = nosurf.Token(r)

	return td, nil
}

func (m *Templates) Render(rw http.ResponseWriter, r *http.Request, tmpl string, td *model.Page) error {
	if rw == nil {
		return errors.New("responseWriter is nil")
	}

	if r == nil {
		return errors.New("request is nil")
	}

	if td == nil {
		td = new(model.Page)
	}

	// cache if in development mode
	if m.inProd {
		templates, err := parseTemplates(m.basePath)
		if err != nil {
			return err
		}
		m.cache = templates
	}

	// get requested template
	rt, ok := m.cache[tmpl+".gohtml"]
	if !ok {
		return fmt.Errorf("%s.gohtml not found", tmpl)
	}

	td, err := m.addDefaultData(td, r)
	if err != nil {
		return err
	}

	// rendering template
	return rt.Execute(rw, td)
}
