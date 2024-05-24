package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort                string
	ServerHost                string
	SecretKey                 string
	AllowedIssuers            []string
	RedisAddr                 string
	RedisPassword             string
	RedisDB                   int
	AuthAPIEndpoint           string
	DriverLocationApiEndpoint string
}

func LoadConfig(env string) *Config {
	envFile := fmt.Sprintf("%s.env", env)
	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("No .env file found for environment: %s\n", env)
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		redisDB = 0
	}

	return &Config{
		ServerPort:                getEnv("SERVER_PORT", "8081"),
		ServerHost:                getEnv("SERVER_HOST", "localhost"),
		SecretKey:                 getEnv("SECRET_KEY", "valid-token"),
		AllowedIssuers:            strings.Split(getEnv("ALLOWED_ISSUERS", "User-API"), ","),
		RedisAddr:                 getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:             getEnv("REDIS_PASSWORD", ""),
		RedisDB:                   redisDB,
		AuthAPIEndpoint:           getEnv("AUTH_URL", "http://localhost:8082"),
		DriverLocationApiEndpoint: getEnv("DRIVER_LOCATION_API_URL", "http://localhost:8080"),
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
