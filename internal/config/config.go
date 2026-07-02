package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// config
type Config struct {
	Server   ServerConfig
	Database DBConfig
	JWT      JWTConfig
	AWS      AWSConfig
	Upload   UploadConfig
}

// ServerConfig fields
type ServerConfig struct {
	Port    string
	GinMode string
}

// DBConfig fields
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JWT configuration
type JWTConfig struct {
	Secret              string
	ExpiresIn           time.Duration
	RefreshTokenExpires time.Duration
}

// AWS configuration
type AWSConfig struct {
	Region     string
	AccessKey  string
	SecretKey  string
	S3Bucket   string
	S3Endpoint string
}

// Upload config
type UploadConfig struct {
	Path        string
	MaxFileSize int64
}

// load function loads config
func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtExpiresIn, _ := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))
	refreshTokenExpires, _ := time.ParseDuration((getEnv("REFRESH_TOKEN_EXPIRES_IN", "720h")))
	maxUploadSize, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)

	return &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "1357"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
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
			Region:     getEnv("REGION", ""),
			AccessKey:  getEnv("ACCESS_KEY", "test"),
			SecretKey:  getEnv("SECRET_KEY", "test"),
			S3Bucket:   getEnv("S3BUCKET", "bookvault-uploads"),
			S3Endpoint: getEnv("S3ENDPOINT", "http://localhost:4566"),
		},
		Upload: UploadConfig{
			Path:        getEnv("PATH", "./uploads"),
			MaxFileSize: maxUploadSize,
		},
	}, nil

}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
