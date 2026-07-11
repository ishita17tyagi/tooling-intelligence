package server

import (
	"net/http"

	"tooling-intelligence/internal/config"
	"tooling-intelligence/internal/handlers"
)

func NewRouter(
	cfg *config.Config,
	researchHandler *handlers.ResearchResponse,
	batchHandler *handlers.BatchHandler,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.Health(cfg))
	mux.HandleFunc("/health", handlers.Health(cfg))

	mux.HandleFunc("/research", researchHandler.Handle)
	mux.HandleFunc("/batch", batchHandler.Handle)

	return mux
}
