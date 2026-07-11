package report

import (
	"fmt"
	"html/template"
	"math"
	"os"
	"path/filepath"

	"tooling-intelligence/internal/models"
)

type Generator struct{}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(results []models.ResearchResponse, analytics *models.Analytics, outputFile string) error {
	err := os.MkdirAll(filepath.Dir(outputFile), 0755)
	if err != nil {
		return err
	}

	templatePath := filepath.Join("templates", "report.html")
	fmt.Println("========== REPORT ==========")
	fmt.Println("Template:", templatePath)
	fmt.Println("Output:", outputFile)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	fmt.Println("Template parsed")

	report := models.ReportData{
		Analytics: *analytics,
		Results:   results,
		AverageConfidencePercent: int(math.Round(
			analytics.AverageConfidence * 100,
		)),
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Println("Output file created")

	err = tmpl.Execute(file, report)
	if err != nil {
		return err
	}
	fmt.Println("Report generated successfully")
	fmt.Println("============================")

	return nil
}
