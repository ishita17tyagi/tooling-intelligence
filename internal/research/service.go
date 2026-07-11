package research

import (
	"context"
	"encoding/json"
	"fmt"

	"tooling-intelligence/internal/gemini"
	"tooling-intelligence/internal/models"
	"tooling-intelligence/internal/normalize"
	"tooling-intelligence/internal/prompts"
	"tooling-intelligence/internal/verification"
)

type Service struct {
	client gemini.Client
}

func New(client gemini.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) Research(ctx context.Context, application string) (*models.ResearchResponse, error) {
	finalPrompt := fmt.Sprintf(prompts.ResearchPrompt, application)

	response, err := s.client.GenerateContent(ctx, finalPrompt)
	if err != nil {
		return nil, err
	}

	var responseData models.ResearchResponse

	fmt.Println("========== GEMINI ==========")
	fmt.Println(response)
	fmt.Println("============================")

	err = json.Unmarshal([]byte(response), &responseData.Result)
	if err != nil {
		fmt.Println("UNMARSHAL ERROR:", err)
		fmt.Println(response)
		return nil, err
	}

	normalizer := normalizer.New()
	normalizer.Normalize(&responseData.Result)

	verifier := verification.New()
	verificationResult := verifier.Verify(&responseData.Result)

	fmt.Println()
	fmt.Println("========== VERIFICATION ==========")
	fmt.Printf("Structural Score   : %.2f\n", verificationResult.StructuralScore)
	fmt.Printf("Evidence Score     : %.2f\n", verificationResult.EvidenceScore)
	fmt.Printf("Completeness Score : %.2f\n", verificationResult.CompletenessScore)
	fmt.Printf("Overall Confidence : %.2f\n", verificationResult.OverallConfidence)

	if len(verificationResult.Deductions) == 0 {
		fmt.Println("Deductions         : None")
	} else {
		fmt.Println("Deductions:")
		for _, d := range verificationResult.Deductions {
			fmt.Printf(" - %s\n", d)
		}
	}
	fmt.Println("==================================")

	responseData.Verification = *verificationResult

	return &responseData, nil
}
