package main

import (
	"errors"
	"github.com/3stadt/movie-nights/imdb"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/3stadt/movie-nights/db"
	"github.com/labstack/echo/v4"
	"github.com/nfnt/resize"
)

type Handler struct {
	DB      *db.DB
	Lang    string
	ImdbApi imdb.Config
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

func (h *Handler) result(c echo.Context) error {
	term := c.QueryParam("q")
	results, err := h.ImdbApi.SearchMovie(h.Lang, term)
	if err != nil {
		return c.Render(http.StatusOK, "result", struct {
			ErrorMessage string
		}{err.Error()})
	}

	for _, res := range results.Results {
		ext := path.Ext(res.Image)
		imgName := res.MovieID + ext

		_, err = os.Open("static/cache/" + imgName)
		if errors.Is(err, os.ErrNotExist) {
			err = cacheImage(res.Image, imgName)
			if err != nil {
				return c.Render(http.StatusOK, "movie_detail", struct {
					ErrorMessage string
				}{err.Error()})
			}
		}

		res.Image = "/static/cache/" + imgName
	}

	return c.Render(http.StatusOK, "result", results)
}

func (h *Handler) movieDetail(c echo.Context) error {
	movieID := c.Param("id")

	movie, err := h.DB.GetMovieFromCache(movieID)
	if err != nil {
		return c.Render(http.StatusOK, "movie_detail", struct {
			ErrorMessage string
		}{err.Error()})
	}
	if movie == nil {
		movie, err = h.ImdbApi.MovieDetail(h.Lang, movieID)
		if err != nil {
			return c.Render(http.StatusOK, "movie_detail", struct {
				ErrorMessage string
			}{err.Error()})
		}
		h.DB.CacheMovie(movie)
	}

	ext := path.Ext(movie.Image)
	imgName := movieID + ext

	_, err = os.Open("static/cache/" + imgName)
	if errors.Is(err, os.ErrNotExist) {
		err = cacheImage(movie.Image, imgName)
		if err != nil {
			return c.Render(http.StatusOK, "movie_detail", struct {
				ErrorMessage string
			}{err.Error()})
		}
	}

	movie.Image = "/static/cache/" + imgName

	return c.Render(http.StatusOK, "movie_detail", movie)
}

func cacheImage(URL, fileName string) error {

	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}

	file, err := os.Create("static/cache/" + strings.TrimSuffix(fileName, path.Ext(fileName)) + ".jpg")
	if err != nil {
		return err
	}

	defer file.Close()

	src, _, err := image.Decode(response.Body)
	if err != nil {
		return err
	}
	newImage := resize.Resize(250, 0, src, resize.Lanczos3)

	// Encode uses a Writer, use a Buffer if you need the raw []byte
	err = jpeg.Encode(file, newImage, nil)

	return nil
}
