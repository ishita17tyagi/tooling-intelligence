package models

type VerificationResult struct {
	StructuralScore   float64  `json:"structural_score"`
	EvidenceScore     float64  `json:"evidence_score"`
	CompletenessScore float64  `json:"completeness_score"`
	OverallConfidence float64  `json:"overall_confidence"`
	Deductions        []string `json:"deductions"`
}
