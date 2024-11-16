package handlers

import (
	"encoding/json"
	"github.com/alekseiapa/go-gl-redundant-pipeline-cleaner/internal/config"
	"github.com/alekseiapa/go-gl-redundant-pipeline-cleaner/internal/gitlab"
	"log"
	"net/http"
)

type WebhookHandler struct {
	Config       *config.Config
	GitlabClient *gitlab.GitlabClient
	Logger       *log.Logger
}

func NewWebhookHandler(cfg *config.Config, glClient *gitlab.GitlabClient, logger *log.Logger) *WebhookHandler {
	return &WebhookHandler{
		Config:       cfg,
		GitlabClient: glClient,
		Logger:       logger,
	}
}

func (h *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.Logger.Printf("Invalid Payload Structure: %v", err)
		http.Error(w, "Invalid Payload Structure", http.StatusBadRequest)
		return
	}
	mrDetails, ok := payload["object_attributes"].(map[string]interface{})
	if !ok {
		h.Logger.Println("Invalid payload structure")
		http.Error(w, "Invalid Payload Structure", http.StatusBadRequest)
		return
	}

	mrID := int(mrDetails["iid"].(float64))
	action := mrDetails["action"].(string)

	log.Printf("Received webhook for MR ID %v with action %v", mrID, action)

	go func() {
		err := h.GitlabClient.CancelRedundantPipelinesByMR(mrID, action)
		if err != nil {
			h.Logger.Fatal("Failed to cancel redundant pipelines for MR ID %v: %v", mrID, err)
		}
	}()

	w.WriteHeader(http.StatusAccepted)
}
