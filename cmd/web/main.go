package main

import (
	"bytes"
	"context"
	"fmt"
	"go-boilerplate/config"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	S3Config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	S3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

// main initializes the application by loading configuration, setting up logging,
// database connection, validation, and fiber app. It also configures error handlers
// and bootstraps the application components before starting the server on the
// configured port.

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validator := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)
	s3 := config.NewS3(log, viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Validate: validator,
		Log:      log,
		Config:   viperConfig,
		S3:       s3,
	})

	app.Post("/test", func(c *fiber.Ctx) error {

		var (
			accessKey = viperConfig.GetString("s3.access_key")
			secretKey = viperConfig.GetString("s3.secret_key")
			endpoint  = viperConfig.GetString("s3.endpoint")
			region    = viperConfig.GetString("s3.region")
		)

		cfg, err := S3Config.LoadDefaultConfig(context.TODO())
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}

		client := S3.NewFromConfig(cfg, func(o *S3.Options) {
			o.Region = region
			o.UsePathStyle = true
			o.EndpointOptions.DisableHTTPS = true
			o.BaseEndpoint = aws.String(endpoint)
			o.Credentials = credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
		})

		uploader := manager.NewUploader(client)

		file, err := c.FormFile("avatar")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid file")
		}
		fileBytes, _ := file.Open()

		defer fileBytes.Close()
		path := "avatar/" + file.Filename
		fileContent, _ := io.ReadAll(fileBytes)
		result, err := uploader.Upload(context.TODO(), &S3.PutObjectInput{
			Bucket: aws.String("go-boilerplate"),
			Key:    aws.String(path),
			Body:   bytes.NewReader(fileContent),
		})

		if err != nil {
			log.Errorf("Failed to upload file : %+v", err.Error())
			return err
		}

		return c.JSON(fiber.Map{
			"message": result,
		})
	})

	webPort := viperConfig.GetInt32("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
