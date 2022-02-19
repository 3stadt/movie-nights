package models

import "gorm.io/gorm"

type ImdbMovie struct {
	*gorm.Model
	MovieID string `gorm:"unique;not null;size:32"`
	JSON    string `gorm:"not null"`
}
