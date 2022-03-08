package models

import (
	"gorm.io/gorm"
	"time"
)

type MovieNight struct {
	*gorm.Model
	Date   time.Time
	Topic  string   `gorm:"unique;not null;size:128"`
	Movies []*Movie `gorm:"many2many:movie_genres;"`
}
