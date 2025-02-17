package repository

import (
	"context"
	"go-boilerplate/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Log *logrus.Logger
	Db  *gorm.DB
}

type UserRepositoryInterface interface {
	Create(c context.Context, entity *entity.User) error
	UpdateProfile(c context.Context, entity *entity.User) (*entity.User, error)
	Show(c context.Context, id string) (*entity.User, error)
	IsEmail(c context.Context, email string) (bool, error)
	IsUsername(c context.Context, username string) (bool, error)
	FindByEmail(c context.Context, email string) (*entity.User, error)
	ResetPassword(c context.Context, entity *entity.User, password string) (*entity.User, error)
	UpdateAvatar(c context.Context, entity *entity.Avatar) (*entity.Avatar, error)
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

func (r *UserRepository) Show(c context.Context, id string) (*entity.User, error) {
	user := new(entity.User)

	err := r.Db.WithContext(c).Preload("Avatar").Model(user).Where("id = ?", id).First(user).Error

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateProfile(c context.Context, entity *entity.User) (*entity.User, error) {
	err := r.Db.WithContext(c).Model(entity).Omit("username").Updates(entity).Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *UserRepository) ResetPassword(c context.Context, entity *entity.User, password string) (*entity.User, error) {
	err := r.Db.WithContext(c).Model(entity).Update("password", password).Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *UserRepository) UpdateAvatar(c context.Context, entity *entity.Avatar) (*entity.Avatar, error) {

	err := r.Db.WithContext(c).Model(entity).Where("user_id = ?", entity.UserID).First(entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err := r.Db.WithContext(c).Model(entity).Create(entity).Error
			if err != nil {
				return nil, err
			}
			return entity, nil
		}
		return nil, err
	}

	err = r.Db.WithContext(c).Model(entity).Updates(entity).Error
	if err != nil {
		return nil, err
	}

	return entity, nil
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
