package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hannasdev/kanban-reports/internal/config"
	"github.com/hannasdev/kanban-reports/internal/menu"
	"github.com/hannasdev/kanban-reports/internal/metrics"
	"github.com/hannasdev/kanban-reports/internal/parser"
	"github.com/hannasdev/kanban-reports/internal/reports"
)

func main() {
	var cfg *config.Config
	var err error
	
	// Parse initial configuration
	cfg, err = config.ParseFlags()
	if err != nil {
		// Enhanced error output with helpful suggestions
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
	
	// Check if interactive mode was requested
	if cfg.Interactive {
		fmt.Println("🎯 Starting Interactive Mode...")
		menuSystem := menu.NewMenu()
		cfg, err = menuSystem.Run()
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
		
		// Show configuration summary
		menuSystem.ShowSummary(cfg)
	} else {
		// CLI mode - show what we're doing
		fmt.Printf("🔄 Kanban Reports - CLI Mode\n")
		fmt.Printf("============================\n")
		showConfigSummary(cfg)
	}

	// Parse CSV file
	fmt.Printf("\n📁 Loading kanban data from: %s\n", cfg.CSVPath)
	csvParser := parser.NewCSVParser(cfg.CSVPath)
	
	// Set delimiter from config
	csvParser.WithDelimiter(cfg.Delimiter)
	
	items, err := csvParser.Parse()
	if err != nil {
		fmt.Printf("❌ Error parsing CSV: %v\n", err)
		fmt.Printf("\n💡 Troubleshooting tips:\n")
		fmt.Printf("   • Check that the file exists and is readable\n")
		fmt.Printf("   • Ensure required columns are present: id, name, estimate, is_completed, completed_at\n")
		fmt.Printf("   • Try different delimiter with --delimiter option\n")
		fmt.Printf("   • For help: %s --help\n", os.Args[0])
		os.Exit(1)
	}

	fmt.Printf("✅ Loaded %d kanban items\n", len(items))

	// Generate report or metrics
	fmt.Printf("\n⚙️  Generating output...\n")
	
	var outputContent string
	
	if cfg.IsMetricsReport() {
		// Generate metrics using the metrics package
		metricsGenerator := metrics.NewGenerator(items)
		metricsGenerator.WithAdHocFilter(cfg.AdHocFilter)

		startDate, endDate := cfg.GetDateRange()
		outputContent, err = metricsGenerator.Generate(cfg.MetricsType, cfg.PeriodType, startDate, endDate, cfg.FilterField)
		if err != nil {
			fmt.Printf("❌ Error generating metrics: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Generate regular report using the reports package
		reporter := reports.NewReporter(items)
		reporter.WithAdHocFilter(cfg.AdHocFilter)

		startDate, endDate := cfg.GetDateRange()
		outputContent, err = reporter.GenerateReport(cfg.ReportType, startDate, endDate, cfg.FilterField)
		if err != nil {
			fmt.Printf("❌ Error generating report: %v\n", err)
			os.Exit(1)
		}
	}

	// Output report
	if cfg.OutputPath != "" {
		// Save to file
		err = os.WriteFile(cfg.OutputPath, []byte(outputContent), 0644)
		if err != nil {
			fmt.Printf("❌ Error writing output to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Output saved to: %s\n", cfg.OutputPath)
		
		// Also show a preview in console
		fmt.Printf("\n📋 Preview (first 500 characters):\n")
		fmt.Printf("%s\n", strings.Repeat("=", 50))
		preview := outputContent
		if len(preview) > 500 {
			preview = preview[:500] + "...\n\n[Full report saved to file]"
		}
		fmt.Printf("%s\n", preview)
	} else {
		// Print to console
		fmt.Printf("\n%s\n", strings.Repeat("=", 60))
		fmt.Printf("📊 RESULTS\n")
		fmt.Printf("%s\n", strings.Repeat("=", 60))
		fmt.Printf("%s\n", outputContent)
		
		// Show helpful next steps
		fmt.Printf("\n💡 Next steps:\n")
		fmt.Printf("   • Save to file: add --output filename.txt\n")
		fmt.Printf("   • Try different time periods: --last 7, --last 30, --last 90\n")
		fmt.Printf("   • Explore other report types: %s --examples\n", os.Args[0])
	}
	
	fmt.Printf("\n🎉 Report generation complete!\n")
}

// showConfigSummary displays the current configuration in CLI mode
func showConfigSummary(cfg *config.Config) {
	fmt.Printf("📋 Configuration:\n")
	fmt.Printf("   📁 CSV File: %s\n", cfg.CSVPath)
	
	if cfg.IsMetricsReport() {
		fmt.Printf("   📈 Mode: Metrics (%s)\n", cfg.MetricsType)
		if cfg.MetricsType == metrics.MetricsTypeThroughput || cfg.MetricsType == metrics.MetricsTypeAll {
			fmt.Printf("   ⏰ Period: %s\n", cfg.PeriodType)
		}
	} else {
		fmt.Printf("   📊 Mode: Report (%s)\n", cfg.ReportType)
	}
	
	// Date range
	if cfg.LastNDays > 0 {
		fmt.Printf("   📅 Date Range: Last %d days\n", cfg.LastNDays)
	} else if !cfg.StartDate.IsZero() && !cfg.EndDate.IsZero() {
		fmt.Printf("   📅 Date Range: %s to %s\n", 
			cfg.StartDate.Format("2006-01-02"), 
			cfg.EndDate.Format("2006-01-02"))
	} else {
		fmt.Printf("   📅 Date Range: All time\n")
	}
	
	fmt.Printf("   🔍 Ad-hoc Filter: %s\n", cfg.AdHocFilter)
	fmt.Printf("   🔗 CSV Delimiter: %s\n", cfg.Delimiter.Name)
	
	if cfg.OutputPath != "" {
		fmt.Printf("   💾 Output: %s\n", cfg.OutputPath)
	} else {
		fmt.Printf("   💾 Output: Console\n")
	}
}