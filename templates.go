package main

import (
	"errors"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

type TemplateRegistry struct {
	templates map[string]*template.Template
}

func buildTemplateRegistry() *TemplateRegistry {
	bp := "public/views/"
	templates := make(map[string]*template.Template)
	templates["index"] = template.Must(template.ParseFiles(
		bp+"layouts/base.gohtml",
		bp+"partials/navbar.gohtml",
		bp+"pages/index.gohtml",
	))
	templates["login"] = template.Must(template.ParseFiles(
		bp+"layouts/base.gohtml",
		bp+"partials/navbar.gohtml",
		bp+"pages/login.gohtml",
	))
	return &TemplateRegistry{
		templates: templates,
	}
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	return tmpl.ExecuteTemplate(w, "base", data)
}
