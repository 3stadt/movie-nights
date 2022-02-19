package main

import (
	"github.com/3stadt/movie-nights/imdb"
	"gopkg.in/yaml.v3"
	"log"
	"os"

	"github.com/3stadt/movie-nights/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	ImdbApiKey string `yaml:"imdb_api_key"`
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
	}

	e := echo.New()

	e.Renderer = buildTemplateRegistry()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.Static("/static", "static")
	e.GET("/login", h.login)
	e.GET("/register", h.register)
	e.GET("/result", h.result)
	e.POST("/register", h.doRegister)
	e.GET("/", h.index)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func hash(pwd string) string {

	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return string(hash)
}
