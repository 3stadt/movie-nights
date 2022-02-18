package models

import "gorm.io/gorm"

type Rating struct {
	*gorm.Model
	Score    uint8
	Excluded bool
	UserID   int
	MovieID  int
}
