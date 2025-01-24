package route

import (
	"go-api-fiber/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	UserController *http.UserController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	// c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Get("/user/:id", c.UserController.Current)
}

// func (c *RouteConfig) SetupAuthRoute() {

// }
