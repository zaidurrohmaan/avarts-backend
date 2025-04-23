package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DBUser             string
	DBPass             string
	DBName             string
	DBPort             string
	DBHost             string
	JWTSecret          string
	GoogleClientID     string
	AWSAccessKey       string
	AWSSecretAccessKey string
	AWSRegion          string
	AWSBucketName      string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DBUser = mustGetEnv("DB_USER")
	DBPass = mustGetEnv("DB_PASS")
	DBName = mustGetEnv("DB_NAME")
	DBPort = mustGetEnv("DB_PORT")
	DBHost = mustGetEnv("DB_HOST")
	JWTSecret = mustGetEnv("JWT_SECRET")
	GoogleClientID = mustGetEnv("GOOGLE_CLIENT_ID")
	AWSAccessKey = mustGetEnv("AWS_ACCESS_KEY")
	AWSSecretAccessKey = mustGetEnv("AWS_SECRET_ACCESS_KEY")
	AWSRegion = mustGetEnv("AWS_REGION")
	AWSBucketName = mustGetEnv("AWS_BUCKET_NAME")
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return val
}
