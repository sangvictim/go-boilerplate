package http

import (
	"go-api-fiber/internal/model"
	"go-api-fiber/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log         *logrus.Logger
	userService *service.UserService
}

func NewUserController(log *logrus.Logger, userService *service.UserService) *UserController {
	return &UserController{
		Log:         log,
		userService: userService,
	}
}

func (c *UserController) Current(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")

	request := &model.GetUserRequest{
		ID: userID,
	}

	response, err := c.userService.Current(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}
