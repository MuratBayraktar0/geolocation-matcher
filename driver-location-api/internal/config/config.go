// pkg/config/config.go
package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoDBURI        string
	MongoDBName       string
	MongoDBCollection string
	ServerPort        string
	ServerHost        string
	SecretKey         string
	AllowedIssuers    []string
}

func LoadConfig(env string) *Config {
	envFile := fmt.Sprintf("%s.env", env)
	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("No .env file found for environment: %s\n", env)
	}
	return &Config{
		MongoDBURI:        getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDBName:       getEnv("MONGODB_NAME", "driver_location"),
		MongoDBCollection: getEnv("MONGODB_COLLECTION", "locations"),
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		ServerHost:        getEnv("SERVER_HOST", "localhost"),
		SecretKey:         getEnv("SECRET_KEY", "valid-token"),
		AllowedIssuers:    strings.Split(getEnv("ALLOWED_ISSUERS", "Matching-API,Driver-API"), ","),
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
