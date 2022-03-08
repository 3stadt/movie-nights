package main

import (
	"errors"
	"fmt"
	"github.com/3stadt/movie-nights/db/models"
	"github.com/3stadt/movie-nights/imdb"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"golang.org/x/crypto/bcrypt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"path"
	"strconv"
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
	Users          []models.User
	Watchlist      *models.WatchList
	MovieNights    []models.MovieNight
	MovieNight     *models.MovieNight
	ErrorMessage   string
	SuccessMessage string
	IsLoggedIn     bool
}

func (h *Handler) index(c echo.Context) error {
	sess := h.getSession(c)
	return c.Render(http.StatusOK, "index", SessionData{IsLoggedIn: isLoggedIn(sess, h.DB)})
}

func (h *Handler) login(c echo.Context) error {
	h.getSession(c)
	return c.Render(http.StatusOK, "login", nil)
}

func (h *Handler) doLogin(c echo.Context) error {
	sess := h.getSession(c)
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
	sess.Values["ID"] = user.ID
	sess.Save(c.Request(), c.Response())

	return c.Render(http.StatusOK, "index", SessionData{IsLoggedIn: isLoggedIn(sess, h.DB)})
}

func (h *Handler) doLogout(c echo.Context) error {
	sess := h.getSession(c)
	sess.Values = make(map[interface{}]interface{})
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusFound, "/login")
}

func (h *Handler) admin(c echo.Context) error {
	sess := h.getSession(c)
	if !isLoggedIn(sess, h.DB) {
		return c.Render(http.StatusUnauthorized, "login", SessionData{ErrorMessage: "You need to be logged in to do that.", IsLoggedIn: false})
	}
	if _, ok := sess.Values["ID"]; !ok {
		return c.Redirect(http.StatusFound, "/logout")
	}
	user := h.DB.GetUserByID(sess.Values["ID"].(uint))
	if user.Level < 9000 {
		return c.Render(http.StatusUnauthorized, "index", SessionData{ErrorMessage: "You need to be admin to do that.", IsLoggedIn: false})
	}
	users := h.DB.GetAllUsers()
	return c.Render(http.StatusOK, "admin", SessionData{IsLoggedIn: isLoggedIn(sess, h.DB), Users: users})
}

func (h *Handler) doAdmin(c echo.Context) error {
	sess := h.getSession(c)
	if !isLoggedIn(sess, h.DB) {
		return c.Render(http.StatusUnauthorized, "login", SessionData{ErrorMessage: "You need to be logged in to do that.", IsLoggedIn: false})
	}
	if _, ok := sess.Values["ID"]; !ok {
		return c.Redirect(http.StatusFound, "/logout")
	}
	user := h.DB.GetUserByID(sess.Values["ID"].(uint))
	if user.Level < 9000 {
		return c.Render(http.StatusUnauthorized, "index", SessionData{ErrorMessage: "You need to be admin to do that.", IsLoggedIn: false})
	}
	userID := c.FormValue("id")
	act := c.FormValue("active")
	active := true
	if act == "false" {
		active = false
	}
	u64, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	h.DB.SetUserStatus(uint(u64), !active)
	users := h.DB.GetAllUsers()
	return c.Render(http.StatusOK, "admin", SessionData{IsLoggedIn: isLoggedIn(sess, h.DB), Users: users})
}

func (h *Handler) register(c echo.Context) error {
	sess := h.getSession(c)
	return c.Render(http.StatusOK, "register", SessionData{IsLoggedIn: isLoggedIn(sess, h.DB)})
}

func (h *Handler) doRegister(c echo.Context) error {
	sess := h.getSession(c)
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
	return c.Render(http.StatusOK, "register_done", SessionData{IsLoggedIn: isLoggedIn(sess, h.DB)})
}

func (h *Handler) result(c echo.Context) error {
	sess := h.getSession(c)
	if !isLoggedIn(sess, h.DB) {
		return c.Render(http.StatusUnauthorized, "login", SessionData{ErrorMessage: "You need to be logged in to do that.", IsLoggedIn: false})
	}
	term := c.QueryParam("q")
	results, err := h.ImdbApi.SearchMovie(h.Lang, term)
	if err != nil {
		return c.Render(http.StatusOK, "result", SessionData{ErrorMessage: err.Error(), IsLoggedIn: isLoggedIn(sess, h.DB)})
	}

	for _, res := range results.Results {
		ext := path.Ext(res.Image)
		imgName := res.MovieID + ext

		_, err = os.Open("static/cache/" + imgName)
		if errors.Is(err, os.ErrNotExist) {
			err = cacheImage(res.Image, imgName)
			if err != nil {
				return c.Render(http.StatusOK, "movie_detail", SessionData{ErrorMessage: err.Error(), IsLoggedIn: isLoggedIn(sess, h.DB)})
			}
		}

		res.Image = "/static/cache/" + imgName
	}

	return c.Render(http.StatusOK, "result", SessionData{Results: results, IsLoggedIn: isLoggedIn(sess, h.DB)})
}

func (h *Handler) movieDetail(c echo.Context) error {
	sess := h.getSession(c)
	if !isLoggedIn(sess, h.DB) {
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
			return c.Render(http.StatusOK, "movie_detail", SessionData{ErrorMessage: err.Error(), IsLoggedIn: isLoggedIn(sess, h.DB)})
		}
		h.DB.CacheMovie(movie)
	}

	ext := path.Ext(movie.Image)
	imgName := movieID + ext

	_, err = os.Open("static/cache/" + imgName)
	if errors.Is(err, os.ErrNotExist) {
		err = cacheImage(movie.Image, imgName)
		if err != nil {
			return c.Render(http.StatusOK, "movie_detail", SessionData{ErrorMessage: err.Error(), IsLoggedIn: isLoggedIn(sess, h.DB)})
		}
	}

	movie.Image = "/static/cache/" + imgName

	return c.Render(http.StatusOK, "movie_detail", SessionData{Movie: movie, IsLoggedIn: isLoggedIn(sess, h.DB)})
}

func (h *Handler) watchlist(c echo.Context) error {
	sess := h.getSession(c)
	if !isLoggedIn(sess, h.DB) {
		return c.Render(http.StatusUnauthorized, "login", SessionData{ErrorMessage: "You need to be logged in to do that.", IsLoggedIn: false})
	}

	id, ok := sess.Values["ID"].(uint)
	if !ok {
		return c.Render(http.StatusBadRequest, "index", SessionData{
			ErrorMessage: fmt.Sprintf("Watchlist for this user does not exist: %v (%T)", sess.Values["ID"], sess.Values["ID"]),
			IsLoggedIn:   true,
		})
	}
	wl := h.DB.GetWatchList(id)
	return c.Render(http.StatusOK, "watchlist", SessionData{IsLoggedIn: isLoggedIn(sess, h.DB), Watchlist: wl})
}

func (h *Handler) addToWatchList(c echo.Context) error {
	sess := h.getSession(c)
	if !isLoggedIn(sess, h.DB) {
		return c.Render(http.StatusUnauthorized, "login", SessionData{ErrorMessage: "You need to be logged in to do that.", IsLoggedIn: false})
	}
	MovieID := c.FormValue("movie-id")
	movie, err := h.DB.GetMovieFromCache(MovieID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "watchlist", SessionData{ErrorMessage: err.Error(), IsLoggedIn: true})
	}
	h.DB.AddMovieToWatchList(movie, sess.Values["ID"].(uint))
	return c.Redirect(http.StatusFound, "/watchlist")
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

func (h *Handler) movieNights(c echo.Context) error {
	sess := h.getSession(c)
	if !isLoggedIn(sess, h.DB) {
		return c.Render(http.StatusUnauthorized, "login", SessionData{ErrorMessage: "You need to be logged in to do that.", IsLoggedIn: false})
	}
	mn := h.DB.GetAllMovieNights()
	return c.Render(http.StatusOK, "movie_nights", SessionData{IsLoggedIn: isLoggedIn(sess, h.DB), MovieNights: mn})
}

func (h *Handler) getSession(c echo.Context) *sessions.Session {
	sess, _ := session.Get("movie-nights", c)
	id, ok := sess.Values["ID"]
	if ok {
		u := h.DB.GetUserByID(id.(uint))
		if !u.Active {
			sess.Values = map[interface{}]interface{}{}
		}
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Save(c.Request(), c.Response())
	return sess
}
