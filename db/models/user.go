package models

import (
	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	Email    string `gorm:"unique;not null;size:128"`
	Password string `gorm:"not null"`
	Active   bool
	Level    uint
	Ratings  []*Rating
}
