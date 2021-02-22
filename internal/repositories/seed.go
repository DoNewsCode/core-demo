package repositories

import (
	"math/rand"

	"github.com/DoNewsCode/core/otgorm"
	"github.com/DoNewsCode/skeleton/internal/entities"
	"gorm.io/gorm"
)

func ProvideSeed() []*otgorm.Seed {
	return []*otgorm.Seed{
		{
			Name: "seeding mysql",
			Run: func(db *gorm.DB) error {
				for i := 0; i < 100; i++ {
					db.Create(&entities.User{
						UserName: RandStringRunes(10),
					})
				}
				return nil
			},
		},
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
