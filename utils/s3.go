package utils

import (
	"avarts/config"
	"avarts/constants"
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	s3v2 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
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

	ext := filepath.Ext(fileHeader.Filename)
	fileKey := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	_, err = s3Svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(config.AWSBucketName),
		Key:         aws.String(fileKey),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", config.AWSBucketName, fileKey)
	return url, nil
}

func IsValidImage(file *multipart.File, fileHeader *multipart.FileHeader, maxSize int64) error {
	if fileHeader.Size > maxSize {
		return fmt.Errorf("%s: %d KB", constants.FileSizeExceeded, maxSize / 1024)
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return fmt.Errorf(constants.InvalidImage)
	}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(*file)
	if err != nil {
		return fmt.Errorf(constants.ReadFileFailed)
	}

	_, _, err = image.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return fmt.Errorf(constants.InvalidImage)
	}

	return nil
}

func DeleteS3File(key string) error {
	ctx := context.TODO()

	cfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	s3Client := s3v2.NewFromConfig(cfg)

	_, err = s3Client.DeleteObject(ctx, &s3v2.DeleteObjectInput{
		Bucket: aws.String(config.AWSBucketName),
		Key:    aws.String(key),
	})

	return err
}