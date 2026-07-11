package handlers

import (
	"encoding/json"
	"net/http"

	"tooling-intelligence/internal/config"
	"tooling-intelligence/internal/models"
)

func Health(cfg *config.Config) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		response := models.HealthResponse{
			Status:      "healthy",
			Application: cfg.AppName,
			Version:     "1.0.0",
			Environment: cfg.Environment,
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(response)
	}
}
