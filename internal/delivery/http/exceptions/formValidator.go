package exceptions

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type validationError struct {
	Field   string `json:"field"` // by passing alt name to ReportError like below
	Message string `json:"message"`
}

func ValidatorError(ctx *fiber.Ctx, err error) error {
	validationErrors := make([]validationError, 0)

	for _, err := range err.(validator.ValidationErrors) {
		e := validationError{
			Field:   err.Field(),
			Message: fmt.Sprintf("%s is %s", err.Field(), err.Tag()),
		}
		validationErrors = append(validationErrors, e)

		_, err := json.MarshalIndent(e, "", "  ")
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": "validation error",
		"errors":  validationErrors,
	})

}
