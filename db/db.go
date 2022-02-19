package db

import (
	"encoding/json"
	"github.com/3stadt/movie-nights/db/models"
	"github.com/3stadt/movie-nights/imdb"
	"github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/gorm"
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
	)
	if err != nil {
		return nil, err
	}
	return &DB{conn: db}, nil
}

func (d *DB) AddUser(name, password string) {
	d.conn.Create(&models.User{
		Name:     name,
		Password: password,
	})
}

func (d *DB) GetMovieFromCache(MovieID string) (*imdb.Movie, error) {
	var m models.ImdbMovie
	d.conn.Where(&models.ImdbMovie{MovieID: MovieID}).First(&m)
	if m.JSON == "" {
		return nil, nil
	}
	movie := imdb.Movie{}
	err := json.Unmarshal([]byte(m.JSON), &m)
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
