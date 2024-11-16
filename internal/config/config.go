package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GitlabWebhookSecret string
	GitlabAPIToken      string
	GitlabURL           string
	GitlabProjectID     string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system default environment variable")
	}
	return &Config{
		GitlabWebhookSecret: getEnv("GITLAB_WEBHOOK_SECRET", ""),
		GitlabAPIToken:      getEnv("GITLAB_API_TOKEN", ""),
		GitlabURL:           getEnv("GITLAB_URL", ""),
		GitlabProjectID:     getEnv("GITLAB_PROJECT_ID", ""),
	}
}

// getEnv retrieves the value of an environment variable.
// If the variable is not set, it returns the provided default value.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
