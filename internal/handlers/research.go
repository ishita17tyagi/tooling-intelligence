package handlers

import (
	"encoding/json"
	"net/http"

	"tooling-intelligence/internal/research"
)

type ResearchResponse struct {
	service *research.Service
}

func NewResearchHandler(service *research.Service) *ResearchResponse {
	return &ResearchResponse{
		service: service,
	}
}

func (h *ResearchResponse) Handle(w http.ResponseWriter, r *http.Request) {

	app := r.URL.Query().Get("app")

	if app == "" {
		http.Error(w, "missing query parameter: app", http.StatusBadRequest)
		return
	}

	result, err := h.service.Research(r.Context(), app)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(result)
}
