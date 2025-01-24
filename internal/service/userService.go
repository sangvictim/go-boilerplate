package service

import (
	"context"
	"go-api-fiber/internal/entity"
	"go-api-fiber/internal/model"
	"go-api-fiber/internal/model/dto"
	"go-api-fiber/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	validate       *validator.Validate
	userRepository *repository.UserRepository
}

func NewUserService(DB *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository) *UserService {
	return &UserService{
		DB:             DB,
		Log:            log,
		validate:       validate,
		userRepository: userRepository,
	}
}

func (c *UserService) Current(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var user entity.User
	err := c.userRepository.FindById(tx, &user, request.ID)
	if err != nil {
		return nil, err
	}

	return dto.UserToResponse(&user), nil
}
