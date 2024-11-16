package middleware

import (
	"github.com/alekseiapa/go-gl-redundant-pipeline-cleaner/internal/config"
	"log"
	"net/http"
	"strings"
)

func AuthMiddleware(cfg *config.Config, logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
