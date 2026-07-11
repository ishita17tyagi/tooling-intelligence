package models

type HealthResponse struct {
	Status      string `json:"status"`
	Application string `json:"application"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
}
