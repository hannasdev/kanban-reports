package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourusername/kanban-reports/internal/parser"
	"github.com/yourusername/kanban-reports/internal/reports"
)

func main() {
	// Define command-line flags
	csvPath := flag.String("csv", "", "Path to the kanban CSV file")
	reportType := flag.String("type", "contributor", "Type of report: contributor, epic, product-area, team")
	startDateStr := flag.String("start", "", "Start date (YYYY-MM-DD)")
	endDateStr := flag.String("end", "", "End date (YYYY-MM-DD)")
	lastNDays := flag.Int("last", 0, "Generate report for the last N days")
	outputPath := flag.String("output", "", "Path to save the report (optional)")

	flag.Parse()

	// Validate inputs
	if *csvPath == "" {
		fmt.Println("Error: CSV file path is required")
		flag.Usage()
		os.Exit(1)
	}

	// Parse report type
	var repType reports.ReportType
	switch *reportType {
	case "contributor":
		repType = reports.ReportTypeContributor
	case "epic":
		repType = reports.ReportTypeEpic
	case "product-area":
		repType = reports.ReportTypeProductArea
	case "team":
		repType = reports.ReportTypeTeam
	default:
		fmt.Printf("Error: Unknown report type: %s\n", *reportType)
		flag.Usage()
		os.Exit(1)
	}

	// Parse date range
	startDate := time.Time{}
	endDate := time.Time{}
	var err error

	// If last N days is specified, it takes precedence
	if *lastNDays > 0 {
		endDate = time.Now()
		startDate = endDate.AddDate(0, 0, -*lastNDays)
	} else {
		// Otherwise use explicit start/end dates
		if *startDateStr != "" {
			startDate, err = time.Parse("2006-01-02", *startDateStr)
			if err != nil {
				fmt.Printf("Error parsing start date: %v\n", err)
				os.Exit(1)
			}
		}

		if *endDateStr != "" {
			endDate, err = time.Parse("2006-01-02", *endDateStr)
			if err != nil {
				fmt.Printf("Error parsing end date: %v\n", err)
				os.Exit(1)
			}
			// Set end date to the end of the day
			endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
	}

	// Parse CSV file
	fmt.Println("Loading kanban data from:", *csvPath)
	csvParser := parser.NewCSVParser(*csvPath)
	items, err := csvParser.Parse()
	if err != nil {
		fmt.Printf("Error parsing CSV: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded %d kanban items\n", len(items))

	// Generate report
	fmt.Println("Generating report...")
	reporter := reports.NewReporter(items)
	report, err := reporter.GenerateReport(repType, startDate, endDate)
	if err != nil {
		fmt.Printf("Error generating report: %v\n", err)
		os.Exit(1)
	}

	// Output report
	if *outputPath != "" {
		// Save to file
		err = os.WriteFile(*outputPath, []byte(report), 0644)
		if err != nil {
			fmt.Printf("Error writing report to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Report saved to: %s\n", *outputPath)
	} else {
		// Print to console
		fmt.Println("\nReport:")
		fmt.Println("-------")
		fmt.Println(report)
	}
}