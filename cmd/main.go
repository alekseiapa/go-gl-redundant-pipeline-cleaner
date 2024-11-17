package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alekseiapa/go-gl-redundant-pipeline-cleaner/internal/config"
	"github.com/alekseiapa/go-gl-redundant-pipeline-cleaner/internal/gitlab"
	"github.com/alekseiapa/go-gl-redundant-pipeline-cleaner/internal/handlers"
	"github.com/alekseiapa/go-gl-redundant-pipeline-cleaner/internal/middleware"
)

func main() {
	cfg := config.LoadConfig()

	logger := log.New(os.Stdout, "webhook-listener", log.LstdFlags|log.Lshortfile)
	if cfg.GitlabAPIToken == "" || cfg.GitlabWebhookSecret == "" || cfg.GitlabProjectID == "" {
		logger.Fatal("Missing required environment variables.")
	}
	gitlabClient, err := gitlab.NewGitlabClient(cfg, logger)
	if err != nil {
		logger.Fatalf("failed to initialize gitlab client: %v", err)
	}
	webhookHandler := handlers.NewWebhookHandler(cfg, gitlabClient, logger)
	authMiddleware := middleware.AuthMiddleware(cfg, logger)
	http.Handle("/cancel-redundant-pipelines", authMiddleware(http.HandlerFunc(webhookHandler.HandleWebhook)))

	log.Println("Starting server on port 5001...")
	if err := http.ListenAndServe(":5001", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
