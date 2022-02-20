package main

import (
	"errors"
	"github.com/3stadt/movie-nights/imdb"
	"golang.org/x/crypto/bcrypt"
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

type SessionData struct {
	Email          string
	Password       string
	Results        *imdb.SearchResult
	Movie          *imdb.Movie
	ErrorMessage   string
	SuccessMessage string
	IsLoggedIn     bool
}

func (h *Handler) index(c echo.Context) error {
	sess := getSession(c)
	return c.Render(http.StatusOK, "index", SessionData{IsLoggedIn: isLoggedIn(sess)})
}

func (h *Handler) login(c echo.Context) error {
	getSession(c)
	return c.Render(http.StatusOK, "login", nil)
}

func (h *Handler) doLogin(c echo.Context) error {
	sess := getSession(c)
	mail := c.FormValue("email")
	pass := c.FormValue("password")
	if !validEmail(mail) {
		return c.Render(http.StatusBadRequest, "login", SessionData{
			Email:        mail,
			Password:     pass,
			ErrorMessage: "Invalid Mail Address",
		})
	}

	user := h.DB.GetUserByMail(mail)

	if !user.Active || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)) != nil {
		return c.Render(http.StatusBadRequest, "login", SessionData{
			Email:        mail,
			Password:     pass,
			ErrorMessage: "Invalid user data",
		})
	}
	sess.Values["isLoggedIn"] = true
	sess.Save(c.Request(), c.Response())

	return c.Render(http.StatusOK, "index", SessionData{IsLoggedIn: isLoggedIn(sess)})
}

func (h *Handler) doLogout(c echo.Context) error {
	sess := getSession(c)
	sess.Values = make(map[interface{}]interface{})
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusFound, "/")
}

func (h *Handler) register(c echo.Context) error {
	sess := getSession(c)
	return c.Render(http.StatusOK, "register", SessionData{IsLoggedIn: isLoggedIn(sess)})
}

func (h *Handler) doRegister(c echo.Context) error {
	sess := getSession(c)
	mail := c.FormValue("email")
	pass := c.FormValue("password")
	if !validEmail(mail) {
		return c.Render(http.StatusBadRequest, "register", SessionData{
			Email:        mail,
			Password:     pass,
			ErrorMessage: "Invalid Mail Address",
		})
	}
	if len(pass) < 8 {
		return c.Render(http.StatusBadRequest, "register", SessionData{
			Email:        mail,
			Password:     pass,
			ErrorMessage: "Password must be at least 8 chars long",
		})
	}
	h.DB.AddUser(mail, hash(pass))
	return c.Render(http.StatusOK, "register_done", SessionData{IsLoggedIn: isLoggedIn(sess)})
}

func (h *Handler) result(c echo.Context) error {
	sess := getSession(c)
	if !isLoggedIn(sess) {
		return c.Render(http.StatusUnauthorized, "login", SessionData{ErrorMessage: "You need to be logged in to do that.", IsLoggedIn: false})
	}
	term := c.QueryParam("q")
	results, err := h.ImdbApi.SearchMovie(h.Lang, term)
	if err != nil {
		return c.Render(http.StatusOK, "result", SessionData{ErrorMessage: err.Error(), IsLoggedIn: isLoggedIn(sess)})
	}

	for _, res := range results.Results {
		ext := path.Ext(res.Image)
		imgName := res.MovieID + ext

		_, err = os.Open("static/cache/" + imgName)
		if errors.Is(err, os.ErrNotExist) {
			err = cacheImage(res.Image, imgName)
			if err != nil {
				return c.Render(http.StatusOK, "movie_detail", SessionData{ErrorMessage: err.Error(), IsLoggedIn: isLoggedIn(sess)})
			}
		}

		res.Image = "/static/cache/" + imgName
	}

	return c.Render(http.StatusOK, "result", SessionData{Results: results, IsLoggedIn: isLoggedIn(sess)})
}

func (h *Handler) movieDetail(c echo.Context) error {
	sess := getSession(c)
	if !isLoggedIn(sess) {
		return c.Render(http.StatusUnauthorized, "login", SessionData{ErrorMessage: "You need to be logged in to do that.", IsLoggedIn: false})
	}
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
			return c.Render(http.StatusOK, "movie_detail", SessionData{ErrorMessage: err.Error(), IsLoggedIn: isLoggedIn(sess)})
		}
		h.DB.CacheMovie(movie)
	}

	ext := path.Ext(movie.Image)
	imgName := movieID + ext

	_, err = os.Open("static/cache/" + imgName)
	if errors.Is(err, os.ErrNotExist) {
		err = cacheImage(movie.Image, imgName)
		if err != nil {
			return c.Render(http.StatusOK, "movie_detail", SessionData{ErrorMessage: err.Error(), IsLoggedIn: isLoggedIn(sess)})
		}
	}

	movie.Image = "/static/cache/" + imgName

	return c.Render(http.StatusOK, "movie_detail", SessionData{Movie: movie, IsLoggedIn: isLoggedIn(sess)})
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
