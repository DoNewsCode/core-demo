package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string
	books []Book
}
