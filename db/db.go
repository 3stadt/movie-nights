package db

import (
	"github.com/3stadt/movie-nights/db/models"
	"github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/gorm"
)

type db struct {
	conn *gorm.DB
}

func Open() (*db, error) {

	gdb, err := gorm.Open(sqlite.Open("movie-nights.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = gdb.AutoMigrate(
		models.Genre{},
		models.Movie{},
		models.Provider{},
		models.Rating{},
		models.User{},
	)
	if err != nil {
		return nil, err
	}
	return &db{conn: gdb}, nil
}

func (d *db) AddUser(name, password string) {
	d.conn.Create(models.User{
		Name:     name,
		Password: "",
	})
}
