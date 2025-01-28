package service

import (
	"context"
	"go-api-fiber/internal/entity"
	"go-api-fiber/internal/model"
	"go-api-fiber/internal/model/dto"
	"go-api-fiber/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	UserRepository repository.UserRepositoryInterface
}

type UserServiceInterface interface {
	Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error)
	Login(ctx context.Context, request *model.LoginUserRequest) (*model.LoginUserResponse, error)
}

func NewUserService(Db *gorm.DB, log *logrus.Logger, userRepository *repository.UserRepository) *UserService {
	return &UserService{
		DB:             Db,
		Log:            log,
		UserRepository: userRepository,
	}
}

func (c *UserService) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := &entity.User{
		Username: request.Username,
		Name:     request.Name,
		Email:    request.Email,
		Password: hashPassword(request.Password),
	}

	isUsername, err := c.UserRepository.IsUsername(ctx, user.Username)
	if err != nil {
		c.Log.Warnf("failed to count username: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if isUsername {
		c.Log.Warnf("Username already exist")
		return nil, fiber.NewError(fiber.StatusConflict, "Username already exist")
	}

	isEmail, err := c.UserRepository.IsEmail(ctx, user.Email)
	if err != nil {
		c.Log.Warnf("failed to count email: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	if isEmail {
		c.Log.Warnf("Email already exist")
		return nil, fiber.NewError(fiber.StatusConflict, "Email already exist")
	}

	if err := c.UserRepository.Create(ctx, user); err != nil {
		c.Log.Warnf("failed to create user: %v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return dto.UserToResponse(user), nil
}

func (c *UserService) Login(ctx context.Context, request *model.LoginUserRequest) (*model.LoginUserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)

	user, err := c.UserRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		c.Log.Warnf("Unauthorized : %+v", err.Error())
		return nil, fiber.ErrUnauthorized
	}

	if user.ID == 0 {
		c.Log.Warnf("Unauthorized : %+v", user)
		return nil, fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Unauthorized : %+v", err.Error())
		return nil, fiber.ErrUnauthorized
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err.Error())
		return nil, fiber.ErrInternalServerError
	}

	token := "123456"
	return dto.LoginUserToReponse(user, token), nil

}

func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}
