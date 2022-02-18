package main

import (
	"net/http"

	"github.com/3stadt/movie-nights/db"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	DB *db.DB
}

func (h *Handler) index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "World")
}

func (h *Handler) login(c echo.Context) error {
	return c.Render(http.StatusOK, "login", "World")
}

func (h *Handler) register(c echo.Context) error {
	return c.Render(http.StatusOK, "register", nil)
}

func (h *Handler) doRegister(c echo.Context) error {
	user := c.FormValue("username")
	pass := c.FormValue("password")
	h.DB.AddUser(user, hash(pass))
	return c.Render(http.StatusOK, "register_done", nil)
}
