package utils

import (
	"avarts/config"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadToS3(file multipart.File, fileHeader *multipart.FileHeader, folder string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.AWSRegion),
		Credentials: credentials.NewStaticCredentials(config.AWSAccessKey, config.AWSSecretAccessKey, ""),
	})
	if err != nil {
		log.Fatal("Failed to create session", err)
		return "", err
	}

	s3Svc := s3.New(sess)

	fileKey := fmt.Sprintf("%s/%s", folder, fileHeader.Filename)

	_, err = s3Svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(config.AWSBucketName),
		Key:    aws.String(fileKey),
		Body:   file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", config.AWSBucketName, fileKey)
	return url, nil
}
