package handlers

import (
	"encoding/json"
	"net/http"

	"tooling-intelligence/internal/analytics"
	"tooling-intelligence/internal/batch"
	"tooling-intelligence/internal/report"
	"tooling-intelligence/internal/storage"
)

type BatchHandler struct {
	processor *batch.Processor
}

func NewBatchHandler(processor *batch.Processor) *BatchHandler {
	return &BatchHandler{
		processor: processor,
	}
}

func (h *BatchHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	inputFile := r.URL.Query().Get("file")
	if inputFile == "" {
		inputFile = "data/input/applications.csv"
	}

	responses, err := h.processor.ProcessCSV(r.Context(), inputFile)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	writer := storage.New()
	err = writer.WriteJSON(responses, "data/output/results.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = writer.WriteCSV(responses, "data/output/results.csv")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	analyzer := analytics.New()
	analysisResults := analyzer.Analyze(responses)

	err = writer.WriteAnalytics(*analysisResults, "data/output/analytics.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	generator := report.New()
	err = generator.Generate(responses, analysisResults, "data/output/report.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":           "success",
		"processed":        len(responses),
		"json_output":      "data/output/results.json",
		"csv_output":       "data/output/results.csv",
		"analytics_output": "data/output/analytics.json",
		"report_output":    "data/output/report.html",
	})
}
