package repositories

import (
	"github.com/DoNewsCode/skeleton/internal/entities"
	"github.com/DoNewsCode/std/pkg/otgorm"
	"gorm.io/gorm"
)

func ProvideMigrator() []*otgorm.Migration {
	return []*otgorm.Migration{
		{
			ID: "202010280100",
			Migrate: func(db *gorm.DB) error {
				type User struct {
					gorm.Model
					UserName string
					books    []entities.Book
				}
				return db.AutoMigrate(
					&User{},
				)
			},
			Rollback: func(db *gorm.DB) error {
				type User struct{}
				return db.Migrator().DropTable(&User{})
			},
		},
	}
}
