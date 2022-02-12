package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.gohtml")),
	}

	e := echo.New()

	e.Renderer = t

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", index)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
