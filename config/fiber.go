package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func NewFiber(viper *viper.Viper) *fiber.App {
	App := fiber.New(fiber.Config{
		AppName: viper.GetString("app.name"),
		Prefork: viper.GetBool("web.prefork"),
	})

	return App
}
