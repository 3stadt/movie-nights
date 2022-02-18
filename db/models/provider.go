package models

import "gorm.io/gorm"

type Provider struct {
	*gorm.Model
	Name string `gorm:"unique;not null;size:32"`
}
