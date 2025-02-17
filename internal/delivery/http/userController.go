package http

import (
	"bytes"
	"context"
	"go-boilerplate/internal/delivery/http/exceptions"
	"go-boilerplate/internal/model"
	"go-boilerplate/internal/service"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log         *logrus.Logger
	Validate    *validator.Validate
	userService service.UserServiceInterface
	S3          *s3.Client
}

func NewUserController(log *logrus.Logger, validate *validator.Validate, userService *service.UserService, s3 *s3.Client) *UserController {
	return &UserController{
		Log:         log,
		Validate:    validate,
		userService: userService,
		S3:          s3,
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("validation error: %v", err)
		return exceptions.ValidatorError(ctx, err)
	}

	response, err := c.userService.Create(ctx.Context(), request)
	if err != nil {
		return err
	}
	ctx.SendStatus(fiber.StatusCreated)
	return ctx.JSON(model.ApiResponse[*model.UserResponse]{
		Message: "Register success",
		Data:    response,
	})
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("validation error: %v", err)
		return exceptions.ValidatorError(ctx, err)
	}
	response, err := c.userService.Login(ctx.Context(), request)
	if err != nil {
		return err
	}

	ctx.Status(fiber.StatusOK)
	return ctx.JSON(model.ApiResponse[*model.LoginUserResponse]{
		Message: "Login Success",
		Data:    response,
	})
}

func (c *UserController) getIdByToken(ctx *fiber.Ctx) string {
	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)

	return userID
}
func (c *UserController) Current(ctx *fiber.Ctx) error {
	userID := c.getIdByToken(ctx)

	user, err := c.userService.Show(ctx.Context(), userID)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	ctx.Status(fiber.StatusOK)
	return ctx.JSON(model.ApiResponse[*model.UserResponse]{
		Message: "Current user",
		Data:    user,
	})
}

func (c *UserController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("validation error: %v", err)
		return exceptions.ValidatorError(ctx, err)
	}

	userID := c.getIdByToken(ctx)
	response, err := c.userService.Update(ctx.Context(), request, userID)
	if err != nil {
		return err
	}
	ctx.Status(fiber.StatusOK)
	return ctx.JSON(model.ApiResponse[*model.UserResponse]{
		Message: "Update Success",
		Data:    response,
	})
}

func (c *UserController) ChangeAvatar(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("avatar")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid file")
	}

	// upload to s3
	result, err := c.uploadToS3(file)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to upload file")
	}

	// update avatar
	userID := c.getIdByToken(ctx)

	if err := c.userService.ChangeAvatar(ctx.Context(), userID, file); err != nil {
		return err
	}
	ctx.Status(fiber.StatusOK)
	return ctx.JSON(fiber.Map{
		"message": "Upload Success",
		"data":    result.Key,
	})
}

func (c *UserController) uploadToS3(file *multipart.FileHeader) (*manager.UploadOutput, error) {
	uploader := manager.NewUploader(c.S3)

	fileBytes, _ := file.Open()
	defer fileBytes.Close()

	path := "avatar/" + file.Filename
	fileContent, _ := io.ReadAll(fileBytes)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("go-boilerplate"),
		Key:    aws.String(path),
		Body:   bytes.NewReader(fileContent),
	})
	if err != nil {
		c.Log.Errorf("Failed to upload file : %+v", err.Error())
		return nil, err
	}

	return result, nil
}
