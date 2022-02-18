package models

import (
	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	Name     string `gorm:"unique;not null;size:32"`
	Password string `gorm:"not null"`
	Active   bool
	Ratings  []*Rating
}
