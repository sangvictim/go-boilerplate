package repository

import (
	"context"
	"go-api-fiber/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Log *logrus.Logger
	Db  *gorm.DB
}

type UserRepositoryInterface interface {
	Create(c context.Context, entity *entity.User) error
	IsEmail(c context.Context, email string) (bool, error)
	IsUsername(c context.Context, username string) (bool, error)
	FindByEmail(c context.Context, email string) (*entity.User, error)
}

func NewUserRepository(db *gorm.DB, log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
		Db:  db,
	}
}

func (r *UserRepository) Create(c context.Context, entity *entity.User) error {
	return r.Db.WithContext(c).Create(entity).Error
}

func (r *UserRepository) IsEmail(c context.Context, email string) (bool, error) {
	var total int64
	user := new(entity.User)

	err := r.Db.WithContext(c).Model(user).Where(entity.User{Email: email}).Count(&total).Error

	if total > 0 {
		return true, err
	}
	return false, err
}

func (r *UserRepository) IsUsername(c context.Context, username string) (bool, error) {
	var total int64
	user := new(entity.User)

	err := r.Db.WithContext(c).Model(user).Where(entity.User{Username: username}).Count(&total).Error

	if total > 0 {
		return true, err
	}
	return false, err
}

func (r *UserRepository) FindByEmail(c context.Context, email string) (*entity.User, error) {
	user := new(entity.User)

	if err := r.Db.WithContext(c).Model(user).Where(entity.User{Email: email}).Find(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
