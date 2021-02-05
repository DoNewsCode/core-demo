package repositories

import (
	"context"

	"github.com/DoNewsCode/skeleton/internal/entities"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (r *UserDao) Find(ctx context.Context, id uint) (*entities.User, error) {
	var u entities.User
	err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find user")
	}
	return &u, nil
}
