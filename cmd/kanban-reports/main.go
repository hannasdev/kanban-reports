package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/hannasdev/kanban-reports/internal/parser"
	"github.com/hannasdev/kanban-reports/internal/reports"
)

func main() {
	// Define command-line flags
	csvPath := flag.String("csv", "", "Path to the kanban CSV file")
	reportType := flag.String("type", "contributor", "Type of report: contributor, epic, product-area, team")
	metricsType := flag.String("metrics", "", "Type of metrics: lead-time, throughput, flow, estimation, age, improvement, all")
	periodType := flag.String("period", "month", "Time period for reports: week, month")
	startDateStr := flag.String("start", "", "Start date (YYYY-MM-DD)")
	endDateStr := flag.String("end", "", "End date (YYYY-MM-DD)")
	lastNDays := flag.Int("last", 0, "Generate report for the last N days")
	outputPath := flag.String("output", "", "Path to save the report (optional)")
	delimiterStr := flag.String("delimiter", "auto", "CSV delimiter: comma, tab, semicolon, or auto for automatic detection")
	adHocFilter := flag.String("ad-hoc", "include", "How to handle ad-hoc requests: include, exclude, only")

	flag.Parse()

	// Validate inputs
	if *csvPath == "" {
		fmt.Println("Error: CSV file path is required")
		flag.Usage()
		os.Exit(1)
	}

	// Parse report type and metrics type
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
		if *metricsType == "" { // Only error if metrics type is also not specified
			fmt.Printf("Error: Unknown report type: %s\n", *reportType)
			flag.Usage()
			os.Exit(1)
		}
	}
	
	// Parse metrics type
	var metType reports.MetricsType
	switch *metricsType {
	case "":
		// No metrics specified, using report type
	case "lead-time":
		metType = reports.MetricsTypeLeadTime
	case "throughput":
		metType = reports.MetricsTypeThroughput
	case "flow":
		metType = reports.MetricsTypeFlow
	case "estimation":
		metType = reports.MetricsTypeEstimation
	case "age":
		metType = reports.MetricsTypeAge
	case "improvement":
		metType = reports.MetricsTypeImprovement
	case "all":
		metType = reports.MetricsTypeAll
	default:
		fmt.Printf("Error: Unknown metrics type: %s\n", *metricsType)
		flag.Usage()
		os.Exit(1)
	}
	
	// Parse period type
	var perType reports.PeriodType
	switch *periodType {
	case "week":
		perType = reports.PeriodTypeWeek
	case "month":
		perType = reports.PeriodTypeMonth
	default:
		fmt.Printf("Error: Unknown period type: %s\n", *periodType)
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

	// Parse ad-hoc filter
	var adHocFilterType reports.AdHocFilterType
	switch *adHocFilter {
	case "include":
		adHocFilterType = reports.AdHocFilterInclude
	case "exclude":
		adHocFilterType = reports.AdHocFilterExclude
	case "only":
		adHocFilterType = reports.AdHocFilterOnly
	default:
		fmt.Printf("Error: Unknown ad-hoc filter type: %s\n", *adHocFilter)
		flag.Usage()
		os.Exit(1)
	}

	// Parse CSV file
	fmt.Println("Loading kanban data from:", *csvPath)
	csvParser := parser.NewCSVParser(*csvPath)
	
	// Set delimiter if specified
	switch *delimiterStr {
	case "comma":
		csvParser.WithDelimiter(',')
	case "tab":
		csvParser.WithDelimiter('\t')
	case "semicolon":
		csvParser.WithDelimiter(';')
	case "auto":
		// Auto-detection is the default
	default:
		fmt.Println("Invalid delimiter specified, using auto-detection")
	}
	
	items, err := csvParser.Parse()
	if err != nil {
		fmt.Printf("Error parsing CSV: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded %d kanban items\n", len(items))

	// Generate report or metrics
	fmt.Println("Generating output...")
	reporter := reports.NewReporter(items).WithAdHocFilter(adHocFilterType)
	
	var outputContent string
	
	if *metricsType != "" {
		// Generate metrics
		outputContent, err = reporter.GenerateMetrics(metType, perType, startDate, endDate)
		if err != nil {
			fmt.Printf("Error generating metrics: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Generate regular report
		outputContent, err = reporter.GenerateReport(repType, startDate, endDate)
		if err != nil {
			fmt.Printf("Error generating report: %v\n", err)
			os.Exit(1)
		}
	}

	// Output report
	if *outputPath != "" {
		// Save to file
		err = os.WriteFile(*outputPath, []byte(outputContent), 0644)
		if err != nil {
			fmt.Printf("Error writing output to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Output saved to: %s\n", *outputPath)
	} else {
		// Print to console
		fmt.Println("\nResults:")
		fmt.Println("-------")
		fmt.Println(outputContent)
	}
}