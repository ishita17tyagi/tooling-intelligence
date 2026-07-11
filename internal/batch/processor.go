package batch

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"tooling-intelligence/internal/models"
	"tooling-intelligence/internal/research"
)

type Processor struct {
	service *research.Service
}

func New(service *research.Service) *Processor {
	return &Processor{
		service: service,
	}
}

func (p *Processor) ProcessCSV(ctx context.Context, inputFile string) ([]models.ResearchResponse, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	_, err = reader.Read()
	if err != nil {
		if err == io.EOF {
			return []models.ResearchResponse{}, nil
		}
		return nil, err
	}

	var responses []models.ResearchResponse

	for {
		select {
		case <-ctx.Done():
			return responses, ctx.Err()
		default:
		}

		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) == 0 {
			continue
		}

		application := strings.TrimSpace(record[0])
		if application == "" {
			continue
		}

		fmt.Printf("Researching %s...\n", application)

		res, err := p.service.Research(ctx, application)
		if err != nil {
			fmt.Printf("✗ %s failed: %v\n", application, err)

			failed := models.ResearchResponse{
				Result: models.ResearchResult{
					Application:    application,
					Category:       "Unknown",
					Description:    "Research failed",
					Authentication: "Unknown",
					SelfServe:      "Unknown",
					APISurface:     "Unknown",
					Buildability:   "Unknown",
					MainBlocker:    err.Error(),
					EvidenceURLs:   []string{},
				},
				Verification: models.VerificationResult{
					StructuralScore:   0,
					EvidenceScore:     0,
					CompletenessScore: 0,
					OverallConfidence: 0,
					Deductions: []string{
						err.Error(),
					},
				},
			}

			responses = append(responses, failed)
			continue
		}

		fmt.Printf("✓ %s completed\n", application)
		responses = append(responses, *res)
	}

	fmt.Println("Batch Complete")
	return responses, nil
}
