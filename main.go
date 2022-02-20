package main

import (
	"context"
	"github.com/3stadt/movie-nights/imdb"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"net/mail"
	"os"
	"os/signal"
	"time"

	"github.com/3stadt/movie-nights/db"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	ImdbApiKey   string `yaml:"imdb_api_key"`
	CookieSecret string `yaml:"cookie_secret"`
	Locale       string `yaml:"locale"`
}

func main() {

	yamlBytes, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	cfg := &Config{}
	err = yaml.Unmarshal(yamlBytes, cfg)
	if err != nil {
		log.Fatal(err)
	}

	gdb, err := db.Open()
	if err != nil {
		log.Fatalf(err.Error())
	}
	h := Handler{
		DB:      gdb,
		ImdbApi: imdb.Config{ApiKey: cfg.ImdbApiKey},
		Lang:    cfg.Locale,
	}

	e := echo.New()

	e.Renderer = buildTemplateRegistry()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(cfg.CookieSecret))))

	// Routes
	e.Static("/static", "static")
	e.GET("/login", h.login)
	e.POST("/login", h.doLogin)
	e.GET("/logout", h.doLogout)
	e.GET("/register", h.register)
	e.GET("/result", h.result)
	e.GET("/admin", h.admin)
	e.GET("/movie/:id", h.movieDetail)
	e.POST("/register", h.doRegister)
	e.GET("/", h.index)

	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func hash(pwd string) string {

	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return string(hash)
}

func getSession(c echo.Context) *sessions.Session {
	sess, _ := session.Get("movie-nights", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Save(c.Request(), c.Response())
	return sess
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isLoggedIn(sess *sessions.Session) bool {
	if loggedIn, ok := sess.Values["isLoggedIn"]; ok {
		return loggedIn.(bool)
	}
	return false
}
