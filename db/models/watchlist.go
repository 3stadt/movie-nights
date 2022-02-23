package models

import "gorm.io/gorm"

type WatchList struct {
	*gorm.Model
	UserID uint
	User   User
	Movies []Movie `gorm:"many2many:movie_watchlists;"`
}
