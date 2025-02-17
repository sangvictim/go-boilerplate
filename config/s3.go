package config

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3Config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewS3(log *logrus.Logger, viper *viper.Viper) *s3.Client {
	var (
		accessKey     = viper.GetString("s3.access_key")
		secretKey     = viper.GetString("s3.secret_key")
		session_token = viper.GetString("s3.session_token")
		endpoint      = viper.GetString("s3.endpoint")
		region        = viper.GetString("s3.region")
		path_style    = viper.GetBool("s3.path_style")
		disable_https = viper.GetBool("s3.disable_https")
	)

	cfg, err := s3Config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return nil
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = region
		o.UsePathStyle = path_style
		o.EndpointOptions.DisableHTTPS = disable_https
		o.BaseEndpoint = aws.String(endpoint)
		o.Credentials = credentials.NewStaticCredentialsProvider(accessKey, secretKey, session_token)
	})

	return client
}
