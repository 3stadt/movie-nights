package main

import (
	"github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
	"log"
)

func main() {

	_, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
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
