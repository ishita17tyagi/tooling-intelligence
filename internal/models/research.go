package models

type ResearchResult struct {
	Application string `json:"application"`

	Category string `json:"category"`

	Description string `json:"description"`

	Authentication string `json:"authentication"`

	SelfServe string `json:"self_serve"`

	APISurface string `json:"api_surface"`

	Buildability string `json:"buildability"`

	MainBlocker string `json:"main_blocker"`

	EvidenceURLs []string `json:"evidence_urls"`
}
