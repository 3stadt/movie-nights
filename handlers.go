package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "World")
}

func login(c echo.Context) error {
	return c.Render(http.StatusOK, "login", "World")
}
