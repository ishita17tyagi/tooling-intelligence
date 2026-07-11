package models

type DistributionItem struct {
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type Analytics struct {
	TotalApplications int `json:"total_applications"`

	AverageConfidence float64 `json:"average_confidence"`

	AuthenticationDistribution map[string]DistributionItem `json:"authentication_distribution"`

	CategoryDistribution map[string]DistributionItem `json:"category_distribution"`

	SelfServeDistribution map[string]DistributionItem `json:"self_serve_distribution"`

	BuildabilityDistribution map[string]DistributionItem `json:"buildability_distribution"`

	TopBlockers map[string]int `json:"top_blockers"`

	ManualReviewCount int `json:"manual_review_count"`
}
