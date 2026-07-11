package storage

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"tooling-intelligence/internal/models"
)

type Writer struct{}

func New() *Writer {
	return &Writer{}
}

func (w *Writer) WriteJSON(results []models.ResearchResponse, outputPath string) error {
	err := os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

func (w *Writer) WriteCSV(results []models.ResearchResponse, outputPath string) error {
	err := os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"Application",
		"Category",
		"Authentication",
		"Self Serve",
		"API Surface",
		"Buildability",
		"Main Blocker",
		"Overall Confidence",
	}

	err = writer.Write(header)
	if err != nil {
		return err
	}

	for _, row := range results {
		record := []string{
			row.Result.Application,
			row.Result.Category,
			row.Result.Authentication,
			row.Result.SelfServe,
			row.Result.APISurface,
			row.Result.Buildability,
			row.Result.MainBlocker,
			fmt.Sprintf("%.2f", row.Verification.OverallConfidence),
		}

		err = writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Writer) WriteAnalytics(analytics models.Analytics, outputPath string) error {
	err := os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(analytics)
}
