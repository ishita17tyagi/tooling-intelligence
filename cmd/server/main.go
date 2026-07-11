package main

import (
	"context"
	"fmt"
	"net/http"

	"tooling-intelligence/internal/batch"
	"tooling-intelligence/internal/config"
	"tooling-intelligence/internal/gemini"
	"tooling-intelligence/internal/handlers"
	"tooling-intelligence/internal/logger"
	"tooling-intelligence/internal/research"
	"tooling-intelligence/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New()

	geminiClient, err := gemini.New(context.Background(), cfg.GeminiAPIKey, cfg.Model)
	if err != nil {
		log.Error(err.Error())
		return
	}

	researchService := research.New(geminiClient)
	researchHandler := handlers.NewResearchHandler(researchService)

	batchProcessor := batch.New(researchService)
	batchHandler := handlers.NewBatchHandler(batchProcessor)

	router := server.NewRouter(
		cfg,
		researchHandler,
		batchHandler,
	)

	address := fmt.Sprintf(":%s", cfg.Port)

	log.Info(
		"Starting HTTP Server",
		"port", cfg.Port,
	)

	err = http.ListenAndServe(address, router)
	if err != nil {
		log.Error("server stopped", "error", err)
	}
}
