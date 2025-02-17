package config

import (
	"go-boilerplate/internal/delivery/http"
	"go-boilerplate/internal/delivery/http/route"
	"go-boilerplate/internal/repository"
	"go-boilerplate/internal/service"

	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	S3       *s3.Client
}

func Bootstrap(config *BootstrapConfig) {
	// setup repository
	userRepository := repository.NewUserRepository(config.DB, config.Log)

	// setup service
	userService := service.NewUserService(config.DB, config.Log, config.Config, config.S3, userRepository)

	// setup controller
	userController := http.NewUserController(config.Log, config.Validate, userService, config.S3)
	storageController := http.NewStorageController(config.Log, config.Validate, config.S3)

	routeConfig := route.RouteConfig{
		App:               config.App,
		Viper:             config.Config,
		UserController:    userController,
		StorageController: storageController,
	}

	routeConfig.Setup()
}
