package db

import (
	"github.com/3stadt/movie-nights/db/models"
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
