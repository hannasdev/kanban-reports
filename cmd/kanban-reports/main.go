package main

import (
	"fmt"
	"os"

	"github.com/hannasdev/kanban-reports/internal/config"
	"github.com/hannasdev/kanban-reports/internal/metrics"
	"github.com/hannasdev/kanban-reports/internal/parser"
	"github.com/hannasdev/kanban-reports/internal/reports"
)

func main() {
	// Parse command-line flags
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Parse CSV file
	fmt.Println("Loading kanban data from:", cfg.CSVPath)
	csvParser := parser.NewCSVParser(cfg.CSVPath)
	
	// Set delimiter from config
	csvParser.WithDelimiter(cfg.Delimiter)
	
	items, err := csvParser.Parse()
	if err != nil {
			fmt.Printf("Error parsing CSV: %v\n", err)
			os.Exit(1)
	}

	fmt.Printf("Loaded %d kanban items\n", len(items))

	// Generate report or metrics
	fmt.Println("Generating output...")
	
	var outputContent string
	
	if cfg.IsMetricsReport() {
		// Generate metrics using the metrics package
		metricsGenerator := metrics.NewGenerator(items)
		metricsGenerator.WithAdHocFilter(cfg.AdHocFilter)

		startDate, endDate := cfg.GetDateRange()
		outputContent, err = metricsGenerator.Generate(cfg.MetricsType, cfg.PeriodType, startDate, endDate, cfg.FilterField)
		if err != nil {
				fmt.Printf("Error generating metrics: %v\n", err)
				os.Exit(1)
		}
	} else {
		// Generate regular report using the reports package
		reporter := reports.NewReporter(items)
		reporter.WithAdHocFilter(cfg.AdHocFilter)

		startDate, endDate := cfg.GetDateRange()
		outputContent, err = reporter.GenerateReport(cfg.ReportType, startDate, endDate, cfg.FilterField)
		if err != nil {
			fmt.Printf("Error generating report: %v\n", err)
			os.Exit(1)
		}
	}

	// Output report
	if cfg.OutputPath != "" {
		// Save to file
		err = os.WriteFile(cfg.OutputPath, []byte(outputContent), 0644)
		if err != nil {
			fmt.Printf("Error writing output to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Output saved to: %s\n", cfg.OutputPath)
	} else {
		// Print to console
		fmt.Println("\nResults:")
		fmt.Println("-------")
		fmt.Println(outputContent)
	}
}