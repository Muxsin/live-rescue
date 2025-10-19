package models

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	Title       string
	Description string
	ImagePath   string
}
