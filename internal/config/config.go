package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config contains all application configuration.
type Config struct {
	Environment string
	Server      ServerConfig
	Database    DBConfig
	JWT         JWTConfig
	AWS         AWSConfig
	Upload      UploadConfig
	SMTP        SMTPConfig
}

// ServerConfig contains HTTP server configuration.
type ServerConfig struct {
	Port           string
	GinMode        string
	AllowedOrigins string
}

// DBConfig contains database connection configuration.
type DBConfig struct {
	Environment string
	Host        string
	Port        string
	User        string
	Password    string
	Name        string
	SSLMode     string
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
	SQSEndpoint    string
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

	cfg := &Config{
		Environment: getEnv("APP_ENV", getEnv("ENV", "development")),
		Server: ServerConfig{
			Port:           getEnv("PORT", "1357"),
			GinMode:        getEnv("GIN_MODE", "debug"),
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
		},
		Database: DBConfig{
			Environment: getEnv("APP_ENV", getEnv("ENV", "development")),
			Host:        getEnv("DB_HOST", "localhost"),
			Port:        getEnv("DB_PORT", "5433"),
			User:        getEnv("DB_USER", "user"),
			Password:    getEnv("DB_PASSWORD", ""),
			Name:        getEnv("DB_NAME", "bookvault"),
			SSLMode:     getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:              getEnv("JWT_SECRET", ""),
			ExpiresIn:           jwtExpiresIn,
			RefreshTokenExpires: refreshTokenExpires,
		},
		AWS: AWSConfig{
			Region:         getEnv("AWS_REGION", ""),
			S3Bucket:       getEnv("AWS_S3_BUCKET", ""),
			S3Endpoint:     getEnv("AWS_S3_ENDPOINT", ""),
			SQSEndpoint:    getEnv("AWS_SQS_ENDPOINT", ""),
			EventQueueName: getEnv("AWS_EVENT_QUEUE_NAME", ""),
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
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil

}

// Validate rejects unsafe or incomplete startup configuration.
func (c *Config) Validate() error {
	missing := make([]string, 0, 4)
	if strings.TrimSpace(c.JWT.Secret) == "" {
		missing = append(missing, "JWT_SECRET")
	}
	if strings.TrimSpace(c.Database.Password) == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if c.Environment == "production" && (strings.TrimSpace(c.Server.AllowedOrigins) == "" || c.Server.AllowedOrigins == "*") {
		missing = append(missing, "CORS_ALLOWED_ORIGINS")
	}

	awsMode := strings.TrimSpace(c.AWS.S3Endpoint) == "" && strings.TrimSpace(c.AWS.SQSEndpoint) == ""
	if awsMode {
		if strings.TrimSpace(c.AWS.Region) == "" {
			missing = append(missing, "AWS_REGION")
		}
		if strings.TrimSpace(c.AWS.S3Bucket) == "" {
			missing = append(missing, "AWS_S3_BUCKET")
		}
		if strings.TrimSpace(c.AWS.EventQueueName) == "" {
			missing = append(missing, "AWS_EVENT_QUEUE_NAME")
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required configuration: %s", strings.Join(missing, ", "))
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
