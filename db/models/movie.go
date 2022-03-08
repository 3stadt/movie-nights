package models

import (
	"gorm.io/gorm"
)

type Movie struct {
	*gorm.Model
	Title       string `gorm:"unique;not null;size:128"`
	FSK         string
	ReleaseYear string
	Genres      []*Genre `gorm:"many2many:movie_genres;"`
	ProviderID  int
	Provider    Provider
	Price       float32
	ImdbID      string `gorm:"unique"`
	Ratings     []*Rating
	WatchLists  []*WatchList  `gorm:"many2many:movie_watchlists;"`
	MovieNights []*MovieNight `gorm:"many2many:movie_genres;"`
}
