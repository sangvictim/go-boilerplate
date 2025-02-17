package service

import (
	"context"
	"fmt"
	"go-boilerplate/internal/entity"
	"go-boilerplate/internal/model"
	"go-boilerplate/internal/model/dto"
	"go-boilerplate/internal/repository"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB             *gorm.DB
	Config         *viper.Viper
	Log            *logrus.Logger
	UserRepository repository.UserRepositoryInterface
	S3             *s3.Client
}

type UserServiceInterface interface {
	Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error)
	Show(ctx context.Context, id string) (*model.UserResponse, error)
	Update(ctx context.Context, request *model.UpdateUserRequest, id string) (*model.UserResponse, error)
	ChangeAvatar(ctx context.Context, userID string, avatar *multipart.FileHeader) error
	Login(ctx context.Context, request *model.LoginUserRequest) (*model.LoginUserResponse, error)
}

func NewUserService(Db *gorm.DB, log *logrus.Logger, config *viper.Viper, s3 *s3.Client, userRepository *repository.UserRepository) *UserService {
	return &UserService{
		DB:             Db,
		Log:            log,
		Config:         config,
		S3:             s3,
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

func (c *UserService) Show(ctx context.Context, id string) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	user, err := c.UserRepository.Show(ctx, id)
	if err != nil {
		c.Log.Warnf("Failed to show user : %+v", err.Error())
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	fmt.Printf("user : %+v", user.Avatar)

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err.Error())
		return nil, fiber.ErrInternalServerError
	}

	return dto.UserToResponse(user), nil
}

func (c *UserService) Update(ctx context.Context, request *model.UpdateUserRequest, id string) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := &entity.User{
		ID:    id,
		Name:  request.Name,
		Email: request.Email,
	}

	if request.Password != "" {
		user.Password = hashPassword(request.Password)
	}

	user, err := c.UserRepository.UpdateProfile(ctx, user)
	if err != nil {
		c.Log.Warnf("Failed to show user : %+v", err.Error())
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err.Error())
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

	if user == nil {
		c.Log.Warnf("Unauthorized : %+v", user)
		return nil, fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Unauthorized : %+v", err.Error())
		return nil, fiber.ErrUnauthorized
	}

	exp := time.Now().Add(time.Hour * 24).Unix()
	token, err := c.generateJWTToken(*user, exp)
	if err != nil {
		c.Log.Warnf("Failed to generate access token : %+v", err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err.Error())
		return nil, fiber.ErrInternalServerError
	}

	return dto.LoginUserToReponse(user, token, exp), nil

}

func (c *UserService) ChangeAvatar(ctx context.Context, userID string, avatar *multipart.FileHeader) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user, err := c.UserRepository.Show(ctx, userID)
	if err != nil {
		c.Log.Warnf("Failed to find user : %+v", err.Error())
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}
	// upload to s3
	avatarRequest := &entity.Avatar{
		UserID:       user.ID,
		OriginalName: avatar.Filename,
		Key:          "avatar/" + avatar.Filename,
		Size:         avatar.Size,
		MimeType:     avatar.Header.Get("Content-Type"),
		Visibility:   "PRIVATE",
	}

	result, err := c.UserRepository.UpdateAvatar(ctx, avatarRequest)
	if err != nil {
		c.Log.Warnf("Failed to update avatar : %+v", err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update avatar")
	}
	c.Log.Info(result)

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err.Error())
		return fiber.ErrInternalServerError
	}

	return nil
}

func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func (c *UserService) generateJWTToken(res entity.User, exp int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": c.Config.GetString("app.name"),
		"sub": res.ID,
		"exp": exp,
	})

	tokenString, err := token.SignedString([]byte(c.Config.GetString("jwt.secret")))

	return tokenString, err
}
