package main

import (
	"github.com/3stadt/movie-nights/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

func main() {

	_, err := db.Open()
	if err != nil {
		log.Fatalf(err.Error())
	}

	e := echo.New()

	e.Renderer = buildTemplateRegistry()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/login", login)
	e.GET("/", index)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
