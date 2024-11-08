package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	GitlabWebhookSecret string
	GitlabAPIToken      string
	GitlabURL           string
	GitlabProjectName   string
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
		GitlabProjectName:   getEnv("GITLAB_PROJECT_NAME", ""),
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