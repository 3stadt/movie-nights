package models

import "gorm.io/gorm"

type Genre struct {
	*gorm.Model
	Name   string   `gorm:"unique;not null;size:32"`
	Movies []*Movie `gorm:"many2many:movie_genres;"`
}
