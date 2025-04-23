package config

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client

func InitAWS() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(AWSRegion))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	S3Client = s3.NewFromConfig(cfg)
}
