package main

import (
	"errors"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type TemplateRegistry struct {
	templates map[string]*template.Template
}

func buildTemplateRegistry() *TemplateRegistry {
	bp := "public/views/"
	templates := make(map[string]*template.Template)

	err := filepath.Walk(bp+"pages", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		if info.IsDir() {
			return nil
		}
		templates[strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))] = template.Must(template.ParseFiles(
			bp+"layouts/base.gohtml",
			bp+"partials/navbar.gohtml",
			bp+"pages/"+info.Name(),
		))
		return nil
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
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
