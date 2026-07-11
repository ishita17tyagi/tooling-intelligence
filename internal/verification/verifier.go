package verification

import (
	"net/url"
	"strings"

	"tooling-intelligence/internal/models"
)

type validationResult struct {
	score      float64
	deductions []string
}

type Verifier struct{}

func New() *Verifier {
	return &Verifier{}
}

func (v *Verifier) Verify(result *models.ResearchResult) *models.VerificationResult {
	structural := v.checkStructure(result)
	evidence := v.checkEvidence(result)
	completeness := v.checkCompleteness(result)

	confidence := (structural.score * 0.40) + (evidence.score * 0.35) + (completeness.score * 0.25)

	allDeductions := append([]string{}, structural.deductions...)
	allDeductions = append(allDeductions, evidence.deductions...)
	allDeductions = append(allDeductions, completeness.deductions...)

	return &models.VerificationResult{
		StructuralScore:   structural.score,
		EvidenceScore:     evidence.score,
		CompletenessScore: completeness.score,
		OverallConfidence: confidence,
		Deductions:        allDeductions,
	}
}

func (v *Verifier) checkStructure(result *models.ResearchResult) validationResult {
	score := 1.0
	var deductions []string

	if result.Application == "" {
		score -= 0.20
		deductions = append(deductions, "Empty application field")
	}
	if result.Category == "" {
		score -= 0.20
		deductions = append(deductions, "Empty category field")
	}
	if result.Description == "" {
		score -= 0.20
		deductions = append(deductions, "Empty description field")
	}
	if result.Authentication == "" {
		score -= 0.20
		deductions = append(deductions, "Empty authentication field")
	}
	if len(result.EvidenceURLs) == 0 {
		score -= 0.20
		deductions = append(deductions, "Missing evidence URLs")
	}

	if score < 0.0 {
		score = 0.0
	}

	return validationResult{score: score, deductions: deductions}
}

func (v *Verifier) checkEvidence(result *models.ResearchResult) validationResult {
	score := 1.0
	var deductions []string

	if len(result.EvidenceURLs) == 0 {
		return validationResult{score: 0.0, deductions: []string{"No evidence URLs provided to validate"}}
	}

	for _, link := range result.EvidenceURLs {
		parsed, err := url.ParseRequestURI(link)
		if err != nil {
			score -= 0.10
			deductions = append(deductions, "Invalid URL format: "+link)
			continue
		}

		if parsed.Scheme != "https" {
			score -= 0.10
			deductions = append(deductions, "Insecure URL scheme: "+link)
			continue
		}

		host := strings.ToLower(parsed.Hostname())
		if host == "localhost" || strings.Contains(host, "example.com") {
			score -= 0.10
			deductions = append(deductions, "Rejected testing domain: "+link)
			continue
		}

		domainOk := false
		allowedSubdomains := []string{"developer.", "docs.", "api.", "developers."}
		for _, sub := range allowedSubdomains {
			if strings.HasPrefix(host, sub) {
				domainOk = true
				break
			}
		}

		if !domainOk {
			allowedRoots := []string{"stripe.com", "shopify.dev"}
			for _, root := range allowedRoots {
				if host == root || strings.HasSuffix(host, "."+root) {
					domainOk = true
					break
				}
			}
		}

		if !domainOk {
			score -= 0.10
			deductions = append(deductions, "URL domain missing official subdomains or trusted roots: "+link)
		}
	}

	if score < 0.0 {
		score = 0.0
	}

	return validationResult{score: score, deductions: deductions}
}

func (v *Verifier) checkCompleteness(result *models.ResearchResult) validationResult {
	score := 1.0
	var deductions []string

	fields := map[string]string{
		"API Surface":  result.APISurface,
		"Buildability": result.Buildability,
		"Main Blocker": result.MainBlocker,
	}

	for name, val := range fields {
		if val == "" {
			score -= 0.20
			deductions = append(deductions, "Empty completeness field: "+name)
		} else if strings.EqualFold(val, "Unknown") {
			score -= 0.15
			deductions = append(deductions, "Unknown value for field: "+name)
		}
	}

	if score < 0.0 {
		score = 0.0
	}

	return validationResult{score: score, deductions: deductions}
}
