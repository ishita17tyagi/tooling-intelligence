package normalizer

import (
	"strings"

	"tooling-intelligence/internal/models"
)

type Normalizer struct{}

func New() *Normalizer {
	return &Normalizer{}
}

func (n *Normalizer) Normalize(result *models.ResearchResult) {
	result.Category = n.normalizeCategory(result.Category)
	result.Authentication = n.normalizeAuthentication(result.Authentication)
	result.Buildability = n.normalizeBuildability(result.Buildability)
	result.SelfServe = n.normalizeSelfServe(result.SelfServe)
}

func (n *Normalizer) normalizeCategory(category string) string {
	val := strings.ToLower(category)

	if strings.Contains(val, "github") || strings.Contains(val, "developer") || strings.Contains(val, "version control") || strings.Contains(val, "source code") {
		return "Developer Tools"
	}
	if strings.Contains(category, "productivity") || strings.Contains(category, "workspace") || strings.Contains(category, "knowledge") {
		return "Productivity"
	}
	if strings.Contains(val, "communication") || strings.Contains(val, "collaboration") || strings.Contains(val, "messaging") || strings.Contains(val, "team") {
		return "Communication"
	}
	if strings.Contains(val, "platform") || strings.Contains(val, "code hosting") || strings.Contains(val, "tool") {
		return "Developer Tools"
	}
	if strings.Contains(val, "productivity") || strings.Contains(val, "workspace") || strings.Contains(val, "knowledge") {
		return "Productivity"
	}

	trimmed := strings.TrimSpace(category)
	if trimmed == "" {
		return "Unknown"
	}
	return trimmed
}

func (n *Normalizer) normalizeAuthentication(auth string) string {
	val := strings.ToLower(auth)
	var methods []string

	if strings.Contains(val, "oauth") {
		methods = append(methods, "OAuth")
	}
	if strings.Contains(val, "api key") || strings.Contains(val, "api_key") || strings.Contains(val, "apikey") {
		methods = append(methods, "API Keys")
	}
	if strings.Contains(val, "personal access token") || strings.Contains(val, "pat") {
		methods = append(methods, "PAT")
	}
	if strings.Contains(val, "jwt") {
		methods = append(methods, "JWT")
	}

	if len(methods) > 0 {
		return strings.Join(methods, " + ")
	}

	trimmed := strings.TrimSpace(auth)
	if trimmed == "" || strings.EqualFold(trimmed, "Unknown") {
		return "Unknown"
	}
	return trimmed
}

func (n *Normalizer) normalizeBuildability(buildability string) string {
	val := strings.ToLower(buildability)

	if strings.Contains(val, "excellent") || strings.Contains(val, "high") || strings.Contains(val, "easy") || strings.Contains(val, "highly") {
		return "High"
	}
	if strings.Contains(val, "medium") || strings.Contains(val, "moderate") {
		return "Medium"
	}
	if strings.Contains(val, "low") || strings.Contains(val, "difficult") || strings.Contains(val, "poor") {
		return "Low"
	}

	return "Unknown"
}

func (n *Normalizer) normalizeSelfServe(selfServe string) string {
	val := strings.ToLower(selfServe)

	if strings.Contains(val, "yes") || strings.Contains(val, "free") || strings.Contains(val, "self") {
		return "Self Serve"
	}
	if strings.Contains(val, "paid") || strings.Contains(val, "enterprise") || strings.Contains(val, "contact sales") || strings.Contains(val, "approval") {
		return "Gated"
	}

	return "Unknown"
}
