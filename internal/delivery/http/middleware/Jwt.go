package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func JWTProtected(v *viper.Viper) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(v.GetString("jwt.secret"))},
		ContextKey:   "user",
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"message": "Missing or malformed JWT"})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"message": "Invalid or expired JWT"})
}
