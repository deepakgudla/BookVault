package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config contains all application configuration.
type Config struct {
	Server   ServerConfig
	Database DBConfig
	JWT      JWTConfig
	AWS      AWSConfig
	Upload   UploadConfig
	SMTP     SMTPConfig
}

// ServerConfig contains HTTP server configuration.
type ServerConfig struct {
	Port    string
	GinMode string
}

// DBConfig contains database connection configuration.
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JWTConfig contains JSON Web Token configuration.
type JWTConfig struct {
	Secret              string
	ExpiresIn           time.Duration
	RefreshTokenExpires time.Duration
}

// AWSConfig contains AWS service configuration.
type AWSConfig struct {
	Region         string
	AccessKey      string
	SecretKey      string
	S3Bucket       string
	EventQueueName string
	S3Endpoint     string
}

// SMTPConfig contains SMTP server configuration.
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// UploadConfig contains file upload configuration.
type UploadConfig struct {
	Path           string
	MaxFileSize    int64
	UploadProvider string
}

// Load reads application configuration from the environment.
func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtExpiresIn, _ := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))
	refreshTokenExpires, _ := time.ParseDuration((getEnv("REFRESH_TOKEN_EXPIRES_IN", "720h")))
	maxUploadSize, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "1025"))

	return &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "1357"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5433"),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "bookvault"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:              getEnv("JWT_SECRET", "<key>"),
			ExpiresIn:           jwtExpiresIn,
			RefreshTokenExpires: refreshTokenExpires,
		},
		AWS: AWSConfig{
			Region:         getEnv("AWS_REGION", ""),
			AccessKey:      getEnv("AWS_ACCESS_KEY_ID", "test"),
			SecretKey:      getEnv("AWS_SECRET_ACCESS_KEY", "test"),
			S3Bucket:       getEnv("AWS_S3_BUCKET", "bookvault-uploads"),
			S3Endpoint:     getEnv("AWS_S3_ENDPOINT", "http://localhost:4566"),
			EventQueueName: getEnv("AWS_EVENT_QUEUE_NAME", "bookvault-events"),
		},
		Upload: UploadConfig{
			Path:           getEnv("UPLOAD_PATH", "./uploads"),
			MaxFileSize:    maxUploadSize,
			UploadProvider: getEnv("UPLOAD_PROVIDER", "local"), // change it to S3 later on
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "localhost"),
			Port:     smtpPort,
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@vault.com"),
		},
	}, nil

}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
