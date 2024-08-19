package env

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration from environment variables
type Config struct {
	MongoDBURI               string
	MongoDBDatabase          string
	SuperTokensConnectionURI string
	SuperTokensAPIKey        string
	SuperTokensAppName       string
	SuperTokensAPIDomain     string
	SuperTokensWebsiteDomain string
	RedisURI                 string
	PORT                     string
}

// Load reads the .env file and loads the environment variables
func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		MongoDBURI:               os.Getenv("MONGODB_URI"),
		MongoDBDatabase:          os.Getenv("MONGODB_DATABASE"),
		SuperTokensConnectionURI: os.Getenv("SUPERTOKENS_CONNECTION_URI"),
		SuperTokensAPIKey:        os.Getenv("SUPERTOKENS_API_KEY"),
		SuperTokensAppName:       os.Getenv("SUPERTOKENS_APP_NAME"),
		RedisURI:                 os.Getenv("REDIS_URI"),
		SuperTokensAPIDomain:     os.Getenv("SUPERTOKENS_API_DOMAIN"),
		SuperTokensWebsiteDomain: os.Getenv("SUPERTOKENS_WEBSITE_DOMAIN"),
		PORT:                     os.Getenv("PORT"),
	}, nil
}

// IsDevelopment checks if the current environment is development
func IsDevelopment() bool {
	return strings.ToLower(os.Getenv("GO_ENV")) == "development"
}
