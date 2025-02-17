package http

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type StorageController struct {
	Log      *logrus.Logger
	Validate *validator.Validate
	S3       *s3.Client
}

func NewStorageController(log *logrus.Logger, validate *validator.Validate, s3 *s3.Client) *StorageController {
	return &StorageController{
		Log:      log,
		Validate: validate,
		S3:       s3,
	}
}

func (c *StorageController) UploadFile(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("avatar")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid file")
	}
	filePath := "/tmp/" + file.Filename
	err = ctx.SaveFile(file, filePath)
	if err != nil {
		return err
	}
	url := ctx.BaseURL() + "/api/storage/upload/" + file.Filename
	ctx.Status(fiber.StatusOK)
	return ctx.JSON(fiber.Map{
		"message": "Upload Success",
		"file":    file.Filename,
		"url":     url,
	})
}

func (c *StorageController) GetFile(ctx *fiber.Ctx) error {
	request := ctx.Params("file")
	return ctx.SendFile("/tmp/" + request)
}

func (c *StorageController) StreamCdnFromS3(ctx *fiber.Ctx) error {
	request := strings.TrimPrefix(ctx.Path(), "/api/cdn/")
	result, err := c.S3.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("go-boilerplate"),
		Key:    aws.String(request),
	})
	if err != nil {
		return err
	}
	return ctx.SendStream(result.Body)
}
