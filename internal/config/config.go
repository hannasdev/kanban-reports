package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hannasdev/kanban-reports/internal/metrics"
	"github.com/hannasdev/kanban-reports/internal/models"
	"github.com/hannasdev/kanban-reports/internal/reports"
	"github.com/hannasdev/kanban-reports/internal/validation"
	"github.com/hannasdev/kanban-reports/pkg/types"
)

// Config represents the application configuration
type Config struct {
	// Input file configuration
	CSVPath     string
	Delimiter   models.DelimiterType
	AutoDetect  bool

	// Report/metrics type configuration
	ReportType  reports.ReportType
	MetricsType metrics.MetricsType
	PeriodType  metrics.PeriodType

	// Date range configuration
	StartDate   time.Time
	EndDate     time.Time
	LastNDays   int

	// Output configuration
	OutputPath  string

	// Filtering configuration
	AdHocFilter types.AdHocFilterType
	FilterField models.FilterField
	
	// CLI mode flags
	Interactive bool
	ShowHelp    bool
}

// ParseFlags parses command-line flags and returns a populated Config
func ParseFlags() (*Config, error) {
	config := &Config{}
	var err error

	// Define command-line flags
	csvPath := flag.String("csv", "", "Path to the kanban CSV file")
	reportType := flag.String("type", "", "Type of report: contributor, epic, product-area, team")
	metricsType := flag.String("metrics", "", "Type of metrics: lead-time, throughput, flow, estimation, age, improvement, all")
	periodType := flag.String("period", "month", "Time period for reports: week, month")
	startDateStr := flag.String("start", "", "Start date (YYYY-MM-DD)")
	endDateStr := flag.String("end", "", "End date (YYYY-MM-DD)")
	lastNDays := flag.Int("last", 0, "Generate report for the last N days")
	outputPath := flag.String("output", "", "Path to save the report (optional)")
	delimiterStr := flag.String("delimiter", "auto", "CSV delimiter: comma, tab, semicolon, or auto for automatic detection")
	adHocFilter := flag.String("ad-hoc", "include", "How to handle ad-hoc requests: include, exclude, only")
	filterField := flag.String("filter-field", "completed_at", "Date field to filter by: completed_at, created_at, started_at")
	
	// New flags for improved CLI experience
	help := flag.Bool("help", false, "Show help information and usage examples")
	helpShort := flag.Bool("h", false, "Show help information and usage examples")
	interactive := flag.Bool("interactive", false, "Run in interactive menu mode")
	interactiveShort := flag.Bool("i", false, "Run in interactive menu mode")
	version := flag.Bool("version", false, "Show version information")
	examples := flag.Bool("examples", false, "Show usage examples")

	// Custom usage function
	flag.Usage = func() {
		showUsage()
	}

	flag.Parse()

	// Handle special flags first
	if *help || *helpShort {
		showHelp()
		os.Exit(0)
	}

	if *version {
		showVersion()
		os.Exit(0)
	}

	if *examples {
		showExamples()
		os.Exit(0)
	}

	if *interactive || *interactiveShort {
		config.Interactive = true
		return config, nil
	}

	// Validate and set CSV path
	if *csvPath == "" {
		return nil, fmt.Errorf("CSV file path is required. Use --csv to specify the file path.\n\nFor help: %s --help", os.Args[0])
	}

	// Validate the CSV path early
	if err := validation.ValidateCSVPath(*csvPath); err != nil {
		csvErr, ok := err.(validation.CSVPathError)
		if !ok {
			return nil, fmt.Errorf("CSV file validation failed: %v\n\nFor help: %s --help", err, os.Args[0])
		}
		
		// Provide specific error messages based on error type
		switch csvErr.Type {
		case "is_directory":
			suggestions := validation.SuggestCSVFiles(*csvPath)
			if len(suggestions) > 0 {
				return nil, fmt.Errorf("%s\n\nFound CSV files in that directory:\n%s\n\nPlease specify the full path to one of these files.\nFor help: %s --help", 
					csvErr.Message, 
					formatSuggestions(suggestions), 
					os.Args[0])
			}
			return nil, fmt.Errorf("%s\n\nFor help: %s --help", csvErr.Message, os.Args[0])
			
		case "not_found":
			return nil, fmt.Errorf("%s\n\nMake sure the file path is correct and the file exists.\nFor help: %s --help", csvErr.Message, os.Args[0])
			
		case "not_readable":
			return nil, fmt.Errorf("%s\n\nCheck file permissions or if the file is open in another program.\nFor help: %s --help", csvErr.Message, os.Args[0])
			
		default:
			return nil, fmt.Errorf("%s\n\nFor help: %s --help", csvErr.Message, os.Args[0])
		}
	}

	config.CSVPath = *csvPath

	// Set delimiter
	config.Delimiter, err = models.ParseDelimiter(*delimiterStr)
	if err != nil {
		fmt.Printf("Warning: %v, using auto-detection\n", err)
		config.Delimiter = models.DelimiterAuto
	}

	// Validate that either report type or metrics type is specified
	if *metricsType == "" && *reportType == "" {
		return nil, fmt.Errorf("either --type or --metrics must be specified.\n\nFor help: %s --help\nFor examples: %s --examples", os.Args[0], os.Args[0])
	}

	// Parse report type with validation
	if *metricsType == "" {
		config.ReportType, err = reports.ParseReportType(*reportType)
		if err != nil {
			return nil, fmt.Errorf("%v\n\nAvailable report types: contributor, epic, product-area, team\nFor help: %s --help", err, os.Args[0])
		}
	}

	// Parse metrics type with validation
	if *metricsType != "" {
		config.MetricsType, err = metrics.ParseMetricsType(*metricsType)
		if err != nil {
			return nil, fmt.Errorf("%v\n\nAvailable metrics types: lead-time, throughput, flow, estimation, age, improvement, all\nFor help: %s --help", err, os.Args[0])
		}
	}

	// Parse period type with validation
	config.PeriodType, err = metrics.ParsePeriodType(*periodType)
	if err != nil {
		return nil, fmt.Errorf("%v\nFor help: %s --help", err, os.Args[0])
	}

	// Parse ad-hoc filter with validation
	config.AdHocFilter, err = types.ParseAdHocFilterType(*adHocFilter)
	if err != nil {
		return nil, fmt.Errorf("%v\nFor help: %s --help", err, os.Args[0])
	}

	// Validate filter field
	config.FilterField, err = models.ParseFilterField(*filterField)
	if err != nil {
		return nil, fmt.Errorf("%v\nFor help: %s --help", err, os.Args[0])
	}

	// Validate and process date range
	if *lastNDays < 0 {
		return nil, fmt.Errorf("last N days must be a positive number, got: %d\nFor help: %s --help", *lastNDays, os.Args[0])
	}

	// If last N days is specified, it takes precedence
	if *lastNDays > 0 {
		config.LastNDays = *lastNDays
		config.EndDate = time.Now()
		config.StartDate = config.EndDate.AddDate(0, 0, -*lastNDays)
	} else {
		// Otherwise use explicit start/end dates
		if *startDateStr != "" {
			config.StartDate, err = time.Parse("2006-01-02", *startDateStr)
			if err != nil {
				return nil, fmt.Errorf("error parsing start date: %v\nExpected format: YYYY-MM-DD\nFor help: %s --help", err, os.Args[0])
			}
		}

		if *endDateStr != "" {
			config.EndDate, err = time.Parse("2006-01-02", *endDateStr)
			if err != nil {
				return nil, fmt.Errorf("error parsing end date: %v\nExpected format: YYYY-MM-DD\nFor help: %s --help", err, os.Args[0])
			}
			// Set end date to the end of the day
			config.EndDate = config.EndDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
	}

	// Set output path
	config.OutputPath = *outputPath

	// Validate date range consistency
	if !config.StartDate.IsZero() && !config.EndDate.IsZero() && config.EndDate.Before(config.StartDate) {
		return nil, fmt.Errorf("invalid date range: end date (%s) is before start date (%s)\nFor help: %s --help", 
			config.EndDate.Format("2006-01-02"), config.StartDate.Format("2006-01-02"), os.Args[0])
	}

	return config, nil
}

// IsMetricsReport returns true if a metrics report is requested
func (c *Config) IsMetricsReport() bool {
	return c.MetricsType != ""
}

// GetDateRange returns the configured date range
func (c *Config) GetDateRange() (time.Time, time.Time) {
	return c.StartDate, c.EndDate
}

// showUsage displays basic usage information
func showUsage() {
	fmt.Printf(`%s - Generate reports and metrics from Kanban CSV data

USAGE:
    %s [OPTIONS]
    %s --interactive
    %s --help

For detailed help and examples, use:
    %s --help
    %s --examples

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

// showVersion displays version information
func showVersion() {
	fmt.Printf(`Kanban Reports v1.0.0
A tool for generating productivity reports from Kanban board CSV exports.

Build Information:
  Go Version: %s
  Platform: %s

`, getGoVersion(), getPlatform())
}

// showHelp displays comprehensive help information
func showHelp() {
	fmt.Printf(`🔄 Kanban Reports - Help & Usage Guide
=====================================

DESCRIPTION:
    Generate insightful reports and metrics from your Kanban board CSV exports.
    Track team productivity, analyze flow efficiency, and identify improvement opportunities.

QUICK EXIT:
    Type 'q', 'quit', 'exit', or 'bye' at any prompt to exit gracefully.

USAGE:
    %s [OPTIONS]                    # Command-line mode
    %s --interactive                # Interactive menu mode
    %s --help                       # Show this help

MODES:
    🎯 Interactive Mode (Recommended for beginners):
        %s --interactive
        %s -i
        
        Guided step-by-step menu to configure your report.

    ⚡ Command-line Mode (Great for automation):
        %s --csv data.csv --type contributor --last 7

REQUIRED OPTIONS:
    --csv FILE                      Path to your kanban CSV file
    
    Choose ONE of:
    --type TYPE                     Generate a report (see REPORT TYPES)
    --metrics TYPE                  Generate metrics (see METRICS TYPES)

REPORT TYPES (--type):
    contributor                     Story points by person who completed work
    epic                           Story points by epic/initiative
    product-area                   Story points by product area
    team                           Story points by team

METRICS TYPES (--metrics):
    lead-time                      How long items take from creation to completion
    throughput                     Completion rates over time (items & points)
    flow                          Flow efficiency (active vs waiting time)
    estimation                    Estimation accuracy (estimates vs actual time)
    age                           Age analysis of current incomplete work
    improvement                   Month-over-month improvement trends
    all                           Generate all metrics above

DATE FILTERING:
    --last N                       Include only last N days
    --start YYYY-MM-DD             Start date (inclusive)
    --end YYYY-MM-DD               End date (inclusive)
    
    Examples:
    --last 7                       Last week
    --last 30                      Last month
    --last 90                      Last quarter
    --start 2024-01-01 --end 2024-03-31    Q1 2024

AD-HOC REQUEST FILTERING:
    --ad-hoc include               Include all items (default)
    --ad-hoc exclude               Exclude items labeled 'ad-hoc-request'
    --ad-hoc only                  Only items labeled 'ad-hoc-request'

TIME PERIODS (for metrics):
    --period week                  Group by week (for throughput metrics)
    --period month                 Group by month (default)

OUTPUT OPTIONS:
    --output FILE                  Save report to file
                                  (default: display in console)

CSV OPTIONS:
    --delimiter auto               Auto-detect delimiter (default)
    --delimiter comma              Comma-separated values
    --delimiter semicolon          Semicolon-separated values
    --delimiter tab                Tab-separated values

OTHER OPTIONS:
    --filter-field FIELD           Date field to filter by:
                                  completed_at (default), created_at, started_at
    --help, -h                     Show this help
    --examples                     Show usage examples
    --version                      Show version information
    --interactive, -i              Run interactive mode

CSV FILE FORMAT:
    Your CSV must include these columns:
    • id, name, estimate, is_completed, completed_at
    
    Optional but useful columns:
    • owners, epic, team, product_area, type, labels

GETTING STARTED:
    1. Export your kanban data as CSV
    2. Run: %s --interactive
    3. Follow the guided setup
    4. Analyze your results!

For more examples: %s --examples

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

// showExamples displays practical usage examples
func showExamples() {
	fmt.Printf(`🚀 Kanban Reports - Usage Examples
=================================

BASIC REPORTS:
    # Team productivity this month
    %s --csv kanban-data.csv --type contributor --last 30

    # Epic progress over the quarter
    %s --csv kanban-data.csv --type epic --last 90 --output epic-report.txt

    # Product area breakdown for specific period
    %s --csv kanban-data.csv --type product-area --start 2024-01-01 --end 2024-03-31

METRICS ANALYSIS:
    # Analyze lead times by story point size
    %s --csv kanban-data.csv --metrics lead-time --last 90

    # Weekly throughput trends
    %s --csv kanban-data.csv --metrics throughput --period week --last 180

    # Complete metrics analysis
    %s --csv kanban-data.csv --metrics all --last 90 --output full-analysis.txt

FILTERING EXAMPLES:
    # Exclude ad-hoc work to see planned work only
    %s --csv kanban-data.csv --type team --last 30 --ad-hoc exclude

    # Analyze only urgent/ad-hoc requests
    %s --csv kanban-data.csv --metrics throughput --ad-hoc only --last 60

    # Filter by creation date instead of completion date
    %s --csv kanban-data.csv --type contributor --last 30 --filter-field created_at

ADVANCED WORKFLOWS:
    # Generate monthly reports for stakeholders
    %s --csv kanban-data.csv --type epic --last 30 --output monthly-epic-report.txt
    %s --csv kanban-data.csv --metrics improvement --output monthly-trends.txt

    # Analyze different time periods
    %s --csv kanban-data.csv --metrics lead-time --start 2024-01-01 --end 2024-02-29
    %s --csv kanban-data.csv --metrics lead-time --start 2024-03-01 --end 2024-05-31

    # Compare planned vs ad-hoc work
    %s --csv kanban-data.csv --type contributor --last 30 --ad-hoc exclude --output planned.txt
    %s --csv kanban-data.csv --type contributor --last 30 --ad-hoc only --output adhoc.txt

INTERACTIVE MODE:
    # Best for first-time users or complex configurations
    %s --interactive
    %s -i

AUTOMATION EXAMPLES:
    # Weekly team report script
    #!/bin/bash
    DATE=$(date +%%Y-%%m-%%d)
    %s --csv latest-export.csv --type team --last 7 --output "weekly-report-$DATE.txt"
    
    # Monthly metrics dashboard
    %s --csv kanban-export.csv --metrics all --last 30 --output monthly-metrics.txt

CSV FILE TIPS:
    • Export from Jira, Azure DevOps, Linear, or any kanban tool
    • Ensure dates are in YYYY/MM/DD HH:MM:SS format
    • Use 'ad-hoc-request' label for ad-hoc work filtering
    • Semi-colon separate multiple owners: john@co.com;jane@co.com

COMMON WORKFLOWS:
    1. Sprint Review Prep:
       %s --csv sprint-data.csv --type contributor --last 14

    2. Epic Progress Check:
       %s --csv project-data.csv --type epic --start 2024-01-01

    3. Process Improvement Analysis:
       %s --csv flow-data.csv --metrics flow --last 90

    4. Team Performance Review:
       %s --csv team-data.csv --metrics improvement --last 180

Need help? Run: %s --help

`, 
		// Now provide exactly 24 arguments to match the 24 %s placeholders
		os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], 
		os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], 
		os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], 
		os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

// Helper functions for system information
func getGoVersion() string {
	// In a real implementation, you might want to embed this at build time
	return "1.21+"
}

func getPlatform() string {
	// In a real implementation, you might want to embed this at build time
	return "linux/amd64"
}

// Add this helper function at the end of the file:
func formatSuggestions(suggestions []string) string {
	if len(suggestions) == 0 {
		return ""
	}
	
	result := ""
	for i, suggestion := range suggestions {
		if i >= 3 { // Limit to 3 suggestions in CLI mode
			result += fmt.Sprintf("   ... and %d more", len(suggestions)-3)
			break
		}
		result += fmt.Sprintf("   • %s\n", suggestion)
	}
	return strings.TrimSuffix(result, "\n")
}