package http

import (
	"go-boilerplate/internal/delivery/http/exceptions"
	"go-boilerplate/internal/model"
	"go-boilerplate/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log         *logrus.Logger
	Validate    *validator.Validate
	userService service.UserServiceInterface
}

func NewUserController(log *logrus.Logger, validate *validator.Validate, userService *service.UserService) *UserController {
	return &UserController{
		Log:         log,
		Validate:    validate,
		userService: userService,
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
