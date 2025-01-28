package route

import (
	"go-api-fiber/internal/delivery/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type RouteConfig struct {
	App            *fiber.App
	UserController *http.UserController
}

func (c *RouteConfig) Setup() {
	// setup route group
	route := c.App.Group("/api")

	// setup middleware
	c.App.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	c.App.Use(recover.New())
	c.App.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        60,
		Expiration: 30 * time.Second,
	}))
	c.App.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// setup route
	c.SetupGuestRoute(route)
	c.SetupAuthRoute(route)

}

func (c *RouteConfig) SetupGuestRoute(route fiber.Router) {
	route.Post("/auth/register", c.UserController.Register)
	route.Post("/auth/login", c.UserController.Login)
}

func (c *RouteConfig) SetupAuthRoute(route fiber.Router) {

}
