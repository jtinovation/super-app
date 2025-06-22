package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort         int
	DBHost             string
	DBPort             int
	DBUser             string
	DBPassword         string
	DBName             string
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
	JWTSecretKey       string
	JWTExpirationHours int
	Minio              MinioConfig
}

type MinioConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	UseSSL          bool
	Endpoint        string
	URL             string
}

var AppConfig *Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		ServerPort: getEnvAsInt("APP_PORT", 8080),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnvAsInt("DB_PORT", 3306),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "your_db"),

		RedisAddr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		JWTSecretKey:       getEnv("JWT_SECRET_KEY", "secret"),
		JWTExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 72),

		Minio: MinioConfig{
			AccessKeyID:     getEnv("MINIO_ACCESS_KEY_ID", "minioadmin"),
			SecretAccessKey: getEnv("MINIO_SECRET_ACCESS_KEY", "minioadmin"),
			Bucket:          getEnv("MINIO_BUCKET", "your-bucket"),
			UseSSL:          getEnvAsInt("MINIO_USE_SSL", 1) == 1,
			Endpoint:        getEnv("MINIO_ENDPOINT", "127.0.0.1:9000"),
			URL:             getEnv("MINIO_URL", "http://127.0.0.1:9000"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
