package analytics

import (
	"math"
	"strings"

	"tooling-intelligence/internal/models"
)

type Analyzer struct{}

func New() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) Analyze(results []models.ResearchResponse) *models.Analytics {
	totalApps := len(results)

	if totalApps == 0 {
		return &models.Analytics{
			TotalApplications:          0,
			AverageConfidence:          0,
			CategoryDistribution:       make(map[string]models.DistributionItem),
			AuthenticationDistribution: make(map[string]models.DistributionItem),
			SelfServeDistribution:      make(map[string]models.DistributionItem),
			BuildabilityDistribution:   make(map[string]models.DistributionItem),
			TopBlockers:                make(map[string]int),
			ManualReviewCount:          0,
		}
	}

	var sumConfidence float64
	var manualReviewCount int

	tempCategories := make(map[string]int)
	tempAuthentications := make(map[string]int)
	tempSelfServes := make(map[string]int)
	tempBuildabilities := make(map[string]int)
	blockers := make(map[string]int)

	for _, row := range results {
		sumConfidence += row.Verification.OverallConfidence

		isSystemError := false
		blocker := strings.TrimSpace(row.Result.MainBlocker)
		blockerLower := strings.ToLower(blocker)

		if strings.Contains(blockerLower, "503") ||
			strings.Contains(blockerLower, "429") ||
			strings.Contains(blockerLower, "failed to generate") ||
			strings.Contains(blockerLower, "timeout") ||
			strings.Contains(blockerLower, "connection refused") {
			isSystemError = true
		}

		category := row.Result.Category
		buildability := row.Result.Buildability

		if row.Verification.OverallConfidence < 0.80 ||
			category == "Unknown" ||
			buildability == "Unknown" ||
			isSystemError {

			manualReviewCount++
		}

		cat := strings.TrimSpace(row.Result.Category)
		if cat == "" {
			cat = "Unknown"
		}
		tempCategories[cat]++

		auth := normalizeAuth(row.Result.Authentication)
		tempAuthentications[auth]++

		ss := strings.TrimSpace(row.Result.SelfServe)
		if ss == "" {
			ss = "Unknown"
		}
		tempSelfServes[ss]++

		build := strings.TrimSpace(row.Result.Buildability)
		if build == "" {
			build = "Unknown"
		}
		tempBuildabilities[build]++

		if blocker != "" && !strings.EqualFold(blocker, "None") && !isSystemError {
			blockers[blocker]++
		}
	}

	avgConfidence := sumConfidence / float64(totalApps)

	categories := make(map[string]models.DistributionItem)
	authentications := make(map[string]models.DistributionItem)
	selfServes := make(map[string]models.DistributionItem)
	buildabilities := make(map[string]models.DistributionItem)

	for category, count := range tempCategories {
		percentage := math.Round((float64(count)*100/float64(totalApps))*100) / 100
		categories[category] = models.DistributionItem{
			Count:      count,
			Percentage: percentage,
		}
	}

	for auth, count := range tempAuthentications {
		percentage := math.Round((float64(count)*100/float64(totalApps))*100) / 100
		authentications[auth] = models.DistributionItem{
			Count:      count,
			Percentage: percentage,
		}
	}

	for selfServe, count := range tempSelfServes {
		percentage := math.Round((float64(count)*100/float64(totalApps))*100) / 100
		selfServes[selfServe] = models.DistributionItem{
			Count:      count,
			Percentage: percentage,
		}
	}

	for buildability, count := range tempBuildabilities {
		percentage := math.Round((float64(count)*100/float64(totalApps))*100) / 100
		buildabilities[buildability] = models.DistributionItem{
			Count:      count,
			Percentage: percentage,
		}
	}

	return &models.Analytics{
		TotalApplications:          totalApps,
		AverageConfidence:          math.Round(avgConfidence*100) / 100,
		AuthenticationDistribution: authentications,
		CategoryDistribution:       categories,
		SelfServeDistribution:      selfServes,
		BuildabilityDistribution:   buildabilities,
		TopBlockers:                blockers,
		ManualReviewCount:          manualReviewCount,
	}
}

func normalizeAuth(rawAuth string) string {
	val := strings.ToLower(rawAuth)

	hasOAuth := strings.Contains(val, "oauth")
	hasAPIKey := strings.Contains(val, "api key") || strings.Contains(val, "api_key") || strings.Contains(val, "apikey")
	hasPAT := strings.Contains(val, "pat ") || strings.Contains(val, "personal access token") || strings.Contains(val, "pat") && (len(val) == 3 || strings.Contains(val, "+") || strings.Contains(val, ","))

	if hasAPIKey && hasOAuth {
		return "API Keys + OAuth"
	}
	if hasOAuth && hasPAT {
		return "OAuth + PAT"
	}
	if hasAPIKey && hasPAT {
		return "API Keys + PAT"
	}
	if hasOAuth {
		return "OAuth 2.0"
	}
	if hasAPIKey {
		return "API Keys"
	}
	if hasPAT {
		return "PAT"
	}
	if strings.Contains(val, "basic") {
		return "Basic Auth"
	}

	trimmed := strings.TrimSpace(rawAuth)
	if trimmed == "" || strings.EqualFold(trimmed, "Unknown") {
		return "Unknown"
	}

	return trimmed
}
