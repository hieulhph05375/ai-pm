package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type GoTestEvent struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Test    string    `json:"Test"`
	Elapsed float64   `json:"Elapsed"`
	Output  string    `json:"Output"`
}

type PlaywrightReport struct {
	Stats struct {
		Expected   int `json:"expected"`
		Unexpected int `json:"unexpected"`
		Flaky      int `json:"flaky"`
		Skipped    int `json:"skipped"`
	} `json:"stats"`
}

func main() {
	log.Println("Generating QA Report PDF...")

	// 1. Parse Backend Results
	backendPassed, backendFailed, backendTotal := parseBackendResults("../test-results/backend.json")

	// 2. Parse Frontend Results
	frontendPassed, frontendFailed, frontendTotal := parseFrontendResults("../test-results/frontend.json")

	// 3. Generate PDF
	generatePDF(
		backendPassed, backendFailed, backendTotal,
		frontendPassed, frontendFailed, frontendTotal,
		"qa-report.pdf",
	)

	log.Println("Successfully generated qa-report.pdf")
}

func parseBackendResults(filepath string) (passed, failed, total int) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Printf("Warning: Could not open %s: %v", filepath, err)
		return 0, 0, 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event GoTestEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			continue
		}

		// Only count package-level pass/fail to avoid double counting individual tests
		if event.Test == "" {
			if event.Action == "pass" {
				passed++
				total++
			} else if event.Action == "fail" {
				failed++
				total++
			}
		}
	}
	return passed, failed, total
}

func parseFrontendResults(filepath string) (passed, failed, total int) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		log.Printf("Warning: Could not read %s: %v", filepath, err)
		return 0, 0, 0
	}

	var report PlaywrightReport
	if err := json.Unmarshal(file, &report); err != nil {
		log.Printf("Warning: Could not parse %s: %v", filepath, err)
		return 0, 0, 0
	}

	passed = report.Stats.Expected
	failed = report.Stats.Unexpected
	// Flaky tests are technically passed after retries, but we'll consider them strictly here
	total = passed + failed + report.Stats.Flaky + report.Stats.Skipped

	return passed, failed, total
}

func generatePDF(bPass, bFail, bTot, fPass, fFail, fTot int, outputPath string) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 24)
	pdf.CellFormat(190, 20, "AIVOVAN QA Automation Report", "0", 1, "C", false, 0, "")

	// Timestamp
	pdf.SetFont("Arial", "I", 12)
	pdf.SetTextColor(100, 100, 100)
	timestamp := time.Now().Format("January 2, 2006 at 3:04 PM")
	pdf.CellFormat(190, 10, fmt.Sprintf("Generated on: %s", timestamp), "0", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Summary Section
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(190, 10, "Executive Summary", "B", 1, "L", false, 0, "")
	pdf.Ln(5)

	totalTests := bTot + fTot
	totalPassed := bPass + fPass

	passRate := 0.0
	if totalTests > 0 {
		passRate = (float64(totalPassed) / float64(totalTests)) * 100
	}

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(95, 10, fmt.Sprintf("Total Tests Executed: %d", totalTests), "0", 0, "L", false, 0, "")

	// Color code pass rate
	if passRate == 100 {
		pdf.SetTextColor(0, 150, 0) // Green
	} else if passRate >= 90 {
		pdf.SetTextColor(200, 150, 0) // Orange
	} else {
		pdf.SetTextColor(200, 0, 0) // Red
	}
	pdf.CellFormat(95, 10, fmt.Sprintf("Overall Pass Rate: %.1f%%", passRate), "0", 1, "R", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(10)

	// Table Header
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(190, 10, "Test Suite Breakdown", "0", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(70, 10, "Suite", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Passed", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Failed", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Total", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 12)

	// Backend Row
	pdf.CellFormat(70, 10, "Backend API Tests (Go)", "1", 0, "L", false, 0, "")
	pdf.SetTextColor(0, 150, 0)
	pdf.CellFormat(40, 10, fmt.Sprintf("%d", bPass), "1", 0, "C", false, 0, "")
	pdf.SetTextColor(200, 0, 0)
	pdf.CellFormat(40, 10, fmt.Sprintf("%d", bFail), "1", 0, "C", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(40, 10, fmt.Sprintf("%d", bTot), "1", 1, "C", false, 0, "")

	// Frontend Row
	pdf.CellFormat(70, 10, "Frontend E2E Tests", "1", 0, "L", false, 0, "")
	pdf.SetTextColor(0, 150, 0)
	pdf.CellFormat(40, 10, fmt.Sprintf("%d", fPass), "1", 0, "C", false, 0, "")
	pdf.SetTextColor(200, 0, 0)
	pdf.CellFormat(40, 10, fmt.Sprintf("%d", fFail), "1", 0, "C", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(40, 10, fmt.Sprintf("%d", fTot), "1", 1, "C", false, 0, "")

	pdf.Ln(20)

	// Status Note
	if bFail == 0 && fFail == 0 {
		pdf.SetFont("Arial", "I", 12)
		pdf.SetTextColor(0, 150, 0)
		pdf.CellFormat(190, 10, "All tests passed successfully! The system is stable.", "0", 1, "C", false, 0, "")
	} else {
		pdf.SetFont("Arial", "B", 12)
		pdf.SetTextColor(200, 0, 0)
		pdf.CellFormat(190, 10, "Warning: Some tests failed. Please review the detailed logs.", "0", 1, "C", false, 0, "")
	}

	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		log.Fatalf("Error saving PDF: %v", err)
	}
}
