package route

import (
	"go-boilerplate/internal/delivery/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
)

type RouteConfig struct {
	App               *fiber.App
	Viper             *viper.Viper
	UserController    *http.UserController
	StorageController *http.StorageController
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
	// route.Use(middleware.JWTProtected(c.App, c.Viper))

	// route.Use(middleware.JWT(c.Viper))

	// setup user route
	userRoute := route.Group("/user")
	userRoute.Get("/me", c.UserController.Current)
	userRoute.Patch("/update", c.UserController.Update)
	userRoute.Post("/change-avatar", c.UserController.ChangeAvatar)

	// setup storage route
	storageRoute := route.Group("/storage")
	storageRoute.Post("/upload", c.StorageController.UploadFile)
	storageRoute.Get("/upload/:file", c.StorageController.GetFile)

	route.Get("/cdn/*/:key", c.StorageController.StreamCdnFromS3)
}
