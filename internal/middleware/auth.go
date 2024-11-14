package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/alekseiapa/go-gl-redundant-pipeline-cleaner/internal/config"
)

func AuthMiddleware(cfg *config.Config, logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Restrict to POST requests only
			if r.Method != http.MethodPost {
				logger.Printf("Method not allowed: %s", r.Method)
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}

			token := strings.TrimSpace(r.Header.Get("X-Gitlab-Token"))
			if token == "" {
				logger.Println("Missing X-Gitlab-Token header in the request")
				http.Error(w, "Unauthorized: missing X-Gitlab-Token header", http.StatusUnauthorized)
				return
			}
			if token != cfg.GitlabWebhookSecret {
				logger.Println("Invalid X-Gitlab-Token header in the request")
				http.Error(w, "Unauthorized: invalid X-Gitlab-Token header", http.StatusUnauthorized)
				return
			}
			logger.Printf("Authorized request: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}
