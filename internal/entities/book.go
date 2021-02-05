package entities

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	BookName string
}
