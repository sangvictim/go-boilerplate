package config

import (
	"go-boilerplate/internal/delivery/http"
	"go-boilerplate/internal/delivery/http/route"
	"go-boilerplate/internal/repository"
	"go-boilerplate/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Validate *validator.Validate
	Log      *logrus.Logger
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// setup repository
	userRepository := repository.NewUserRepository(config.DB, config.Log)

	// setup service
	userService := service.NewUserService(config.DB, config.Log, config.Config, userRepository)

	// setup controller
	userController := http.NewUserController(config.Log, config.Validate, userService)

	// setup middleware

	routeConfig := route.RouteConfig{
		App:            config.App,
		Viper:          config.Config,
		UserController: userController,
	}

	routeConfig.Setup()
}
