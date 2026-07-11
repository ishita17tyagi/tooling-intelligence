package models

type ResearchResponse struct {
	Result       ResearchResult     `json:"result"`
	Verification VerificationResult `json:"verification"`
}
