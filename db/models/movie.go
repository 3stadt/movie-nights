package models

import (
	"gorm.io/gorm"
	"time"
)

type Movie struct {
	*gorm.Model
	Name        string `gorm:"unique;not null;size:128"`
	FSK         uint8
	ReleaseYear *time.Time
	Genres      []*Genre `gorm:"many2many:movie_genres;"`
	ProviderID  int
	Provider    Provider
	Price       float32
	Ratings     []*Rating
}
