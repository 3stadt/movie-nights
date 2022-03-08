package db

import (
	"encoding/json"
	"github.com/3stadt/movie-nights/db/models"
	"github.com/3stadt/movie-nights/imdb"
	"github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/gorm"
	"time"
)

type DB struct {
	conn *gorm.DB
}

func Open() (*DB, error) {

	db, err := gorm.Open(sqlite.Open("movie-nights.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(
		models.Genre{},
		models.Movie{},
		models.Provider{},
		models.Rating{},
		models.User{},
		models.ImdbMovie{},
		models.MovieNight{},
	)
	if err != nil {
		return nil, err
	}
	return &DB{conn: db}, nil
}

func (d *DB) AddUser(name, password string) {
	d.conn.Create(&models.User{
		Email:    name,
		Password: password,
		Level:    100, // Default user level
	})
}

func (d *DB) GetAllMovieNights() []models.MovieNight {
	var mn []models.MovieNight
	d.conn.Find(&mn)
	return mn
}

func (d *DB) CreateMovieNight(topic string, date time.Time) {
	mn := models.MovieNight{
		Date:  date,
		Topic: topic,
	}
	d.conn.Create(mn)
}

func (d *DB) AddMovieToMovieNight(movieNightId uint, imdbId string) error {
	movieNight := d.conn.Where(&models.MovieNight{}, movieNightId).Association("Movies")
	return movieNight.Append(&models.Movie{
		ImdbID: imdbId,
	})
}

func (d *DB) GetMovieFromCache(MovieID string) (*imdb.Movie, error) {
	var m models.ImdbMovie
	d.conn.Where(&models.ImdbMovie{MovieID: MovieID}).First(&m)
	if m.JSON == "" {
		return nil, nil
	}
	var movie imdb.Movie
	err := json.Unmarshal([]byte(m.JSON), &movie)
	return &movie, err
}

func (d *DB) CacheMovie(movie *imdb.Movie) {
	mJson, err := json.Marshal(movie)
	if err != nil {
		return
	}
	d.conn.Create(&models.ImdbMovie{
		MovieID: movie.MovieID,
		JSON:    string(mJson),
	})
}

func (d *DB) GetUserByMail(mail string) *models.User {
	var u models.User
	d.conn.Where(&models.User{Email: mail, Active: true}).First(&u)
	return &u
}

func (d *DB) GetUserByID(id uint) *models.User {
	var u models.User
	d.conn.First(&u, id)
	return &u
}

func (d *DB) GetAllUsers() []models.User {
	var u []models.User
	d.conn.Find(&u)
	return u
}

func (d *DB) SetUserStatus(id uint, active bool) {
	d.conn.Model(&models.User{}).Where("id = ?", id).Update("active", active)
}

func (d *DB) GetWatchList(id uint) *models.WatchList {
	wl := models.WatchList{}
	d.conn.Preload("Movies").Preload("Movies.Genres").Where(&models.WatchList{UserID: id}).First(&wl)
	if wl.UserID == 0 {
		d.conn.Create(&models.WatchList{
			UserID: id,
		})
		d.conn.Where(&models.WatchList{UserID: id}).First(&wl)
	}
	return &wl
}

func (d *DB) AddMovieToWatchList(i *imdb.Movie, userID uint) {
	wl := models.WatchList{}
	d.conn.Where(&models.WatchList{UserID: userID}).First(&wl)
	if wl.UserID == 0 {
		d.conn.Create(&models.WatchList{
			UserID: userID,
		})
		d.conn.Where(&models.WatchList{UserID: userID}).First(&wl)
	}
	wl.Movies = append(wl.Movies, models.Movie{
		Title:       i.Title,
		FSK:         i.ContentRating,
		ReleaseYear: i.ReleaseDate,
		ImdbID:      i.MovieID, // TODO ADD GENRES and additional data as price, provider
	})
	d.conn.Model(&wl).Updates(wl)
}
