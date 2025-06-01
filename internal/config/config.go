package config

import (
	"flag"
	"fmt"
	"os"
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

// flagSet holds all parsed command-line flags
type flagSet struct {
	csvPath      *string
	reportType   *string
	metricsType  *string
	periodType   *string
	startDateStr *string
	endDateStr   *string
	lastNDays    *int
	outputPath   *string
	delimiterStr *string
	adHocFilter  *string
	filterField  *string
	
	// Control flags
	help         *bool
	helpShort    *bool
	interactive  *bool
	interactiveShort *bool
	version      *bool
	examples     *bool
}

// ParseFlags parses command-line flags and returns a populated Config
func ParseFlags() (*Config, error) {
	flags := defineFlags()
	
	flag.Usage = showUsage
	flag.Parse()

	// Handle special control flags first
	if err := handleControlFlags(flags); err != nil {
		return nil, err
	}

	// Check for interactive mode
	if *flags.interactive || *flags.interactiveShort {
		return &Config{Interactive: true}, nil
	}

	// Parse and validate configuration
	config, err := buildConfig(flags)
	if err != nil {
		return nil, fmt.Errorf("%v\n\nFor help: %s --help", err, os.Args[0])
	}

	return config, nil
}

// defineFlags sets up all command-line flags
func defineFlags() *flagSet {
	return &flagSet{
		csvPath:      flag.String("csv", "", "Path to the kanban CSV file"),
		reportType:   flag.String("type", "", "Type of report: contributor, epic, product-area, team"),
		metricsType:  flag.String("metrics", "", "Type of metrics: lead-time, throughput, flow, estimation, age, improvement, all"),
		periodType:   flag.String("period", DefaultPeriodType, "Time period for reports: week, month"),
		startDateStr: flag.String("start", "", "Start date (YYYY-MM-DD)"),
		endDateStr:   flag.String("end", "", "End date (YYYY-MM-DD)"),
		lastNDays:    flag.Int("last", 0, "Generate report for the last N days"),
		outputPath:   flag.String("output", "", "Path to save the report (optional)"),
		delimiterStr: flag.String("delimiter", DefaultDelimiter, "CSV delimiter: comma, tab, semicolon, or auto for automatic detection"),
		adHocFilter:  flag.String("ad-hoc", DefaultAdHocFilter, "How to handle ad-hoc requests: include, exclude, only"),
		filterField:  flag.String("filter-field", DefaultFilterField, "Date field to filter by: completed_at, created_at, started_at"),
		
		help:             flag.Bool("help", false, "Show help information and usage examples"),
		helpShort:        flag.Bool("h", false, "Show help information and usage examples"),
		interactive:      flag.Bool("interactive", false, "Run in interactive menu mode"),
		interactiveShort: flag.Bool("i", false, "Run in interactive menu mode"),
		version:          flag.Bool("version", false, "Show version information"),
		examples:         flag.Bool("examples", false, "Show usage examples"),
	}
}

// handleControlFlags processes special flags like help, version, examples
func handleControlFlags(flags *flagSet) error {
	if *flags.help || *flags.helpShort {
		showHelp()
		os.Exit(0)
	}

	if *flags.version {
		showVersion()
		os.Exit(0)
	}

	if *flags.examples {
		showExamples()
		os.Exit(0)
	}

	return nil
}

// buildConfig constructs and validates the configuration from parsed flags
func buildConfig(flags *flagSet) (*Config, error) {
	config := &Config{}

	// Validate and set required fields
	if err := setCSVPath(config, *flags.csvPath); err != nil {
		return nil, err
	}

	if err := setDelimiter(config, *flags.delimiterStr); err != nil {
		return nil, err
	}

	if err := setReportAndMetricsTypes(config, *flags.reportType, *flags.metricsType); err != nil {
		return nil, err
	}

	if err := setPeriodType(config, *flags.periodType); err != nil {
		return nil, err
	}

	if err := setFilterOptions(config, *flags.adHocFilter, *flags.filterField); err != nil {
		return nil, err
	}

	if err := setDateRange(config, *flags.startDateStr, *flags.endDateStr, *flags.lastNDays); err != nil {
		return nil, err
	}

	config.OutputPath = *flags.outputPath

	return config, nil
}

// setCSVPath validates and sets the CSV file path
func setCSVPath(config *Config, csvPath string) error {
	if csvPath == "" {
		return fmt.Errorf("CSV file path is required. Use --csv to specify the file path")
	}

	if err := validation.ValidateCSVPath(csvPath); err != nil {
		return formatCSVValidationError(err, csvPath)
	}

	config.CSVPath = csvPath
	return nil
}

// formatCSVValidationError provides user-friendly error messages for CSV validation failures
func formatCSVValidationError(err error, csvPath string) error {
	csvErr, ok := err.(validation.CSVPathError)
	if !ok {
		return fmt.Errorf("CSV file validation failed: %v", err)
	}

	switch csvErr.Type {
	case "is_directory":
		suggestions := validation.SuggestCSVFiles(csvPath)
		if len(suggestions) > 0 {
			return fmt.Errorf("%s\n\nFound CSV files in that directory:\n%s\n\nPlease specify the full path to one of these files", 
				csvErr.Message, formatSuggestions(suggestions))
		}
		return fmt.Errorf("%s", csvErr.Message)
		
	case "not_found":
		return fmt.Errorf("%s\n\nMake sure the file path is correct and the file exists", csvErr.Message)
		
	case "not_readable":
		return fmt.Errorf("%s\n\nCheck file permissions or if the file is open in another program", csvErr.Message)
		
	default:
		return fmt.Errorf("%s", csvErr.Message)
	}
}

// setDelimiter parses and sets the CSV delimiter
func setDelimiter(config *Config, delimiterStr string) error {
	delimiter, err := models.ParseDelimiter(delimiterStr)
	if err != nil {
		fmt.Printf("Warning: %v, using auto-detection\n", err)
		config.Delimiter = models.DelimiterAuto
		return nil
	}
	config.Delimiter = delimiter
	return nil
}

// setReportAndMetricsTypes validates and sets report/metrics types with proper precedence
func setReportAndMetricsTypes(config *Config, reportType, metricsType string) error {
	if metricsType == "" && reportType == "" {
		return fmt.Errorf("either --type or --metrics must be specified")
	}

	// Metrics type takes precedence when both are specified (original behavior)
	if metricsType != "" {
		mt, err := metrics.ParseMetricsType(metricsType)
		if err != nil {
			return fmt.Errorf("%v\n\nAvailable metrics types: lead-time, throughput, flow, estimation, age, improvement, all", err)
		}
		config.MetricsType = mt
		return nil
	}

	// Only parse report type if metrics type is not specified
	if reportType != "" {
		rt, err := reports.ParseReportType(reportType)
		if err != nil {
			return fmt.Errorf("%v\n\nAvailable report types: contributor, epic, product-area, team", err)
		}
		config.ReportType = rt
		return nil
	}

	return nil
}

// setPeriodType parses and validates the period type
func setPeriodType(config *Config, periodType string) error {
	pt, err := metrics.ParsePeriodType(periodType)
	if err != nil {
		return err
	}
	config.PeriodType = pt
	return nil
}

// setFilterOptions parses and sets filtering configuration
func setFilterOptions(config *Config, adHocFilter, filterField string) error {
	af, err := types.ParseAdHocFilterType(adHocFilter)
	if err != nil {
		return err
	}
	config.AdHocFilter = af

	ff, err := models.ParseFilterField(filterField)
	if err != nil {
		return err
	}
	config.FilterField = ff

	return nil
}

// setDateRange validates and sets the date range configuration
func setDateRange(config *Config, startDateStr, endDateStr string, lastNDays int) error {
	if lastNDays < 0 {
		return fmt.Errorf("last N days must be a positive number, got: %d", lastNDays)
	}

	// Last N days takes precedence
	if lastNDays > 0 {
		config.LastNDays = lastNDays
		config.EndDate = time.Now()
		config.StartDate = config.EndDate.AddDate(0, 0, -lastNDays)
		return nil
	}

	// Parse explicit dates
	if err := parseExplicitDates(config, startDateStr, endDateStr); err != nil {
		return err
	}

	// Validate date range consistency
	if !config.StartDate.IsZero() && !config.EndDate.IsZero() && config.EndDate.Before(config.StartDate) {
		return fmt.Errorf("invalid date range: end date (%s) is before start date (%s)", 
			config.EndDate.Format(DateFormat), config.StartDate.Format(DateFormat))
	}

	return nil
}

// parseExplicitDates parses start and end date strings
func parseExplicitDates(config *Config, startDateStr, endDateStr string) error {
	if startDateStr != "" {
		startDate, err := time.Parse(DateFormat, startDateStr)
		if err != nil {
			return fmt.Errorf("error parsing start date: %v\nExpected format: YYYY-MM-DD", err)
		}
		config.StartDate = startDate
	}

	if endDateStr != "" {
		endDate, err := time.Parse(DateFormat, endDateStr)
		if err != nil {
			return fmt.Errorf("error parsing end date: %v\nExpected format: YYYY-MM-DD", err)
		}
		// Set end date to the end of the day
		config.EndDate = endDate.Add(HoursPerDay*time.Hour + MinutesPerHour*time.Minute + SecondsPerMinute*time.Second)
	}

	return nil
}

// IsMetricsReport returns true if a metrics report is requested
func (c *Config) IsMetricsReport() bool {
	return c.MetricsType != ""
}

// GetDateRange returns the configured date range
func (c *Config) GetDateRange() (time.Time, time.Time) {
	return c.StartDate, c.EndDate
}

// formatSuggestions formats a list of file suggestions for display
func formatSuggestions(suggestions []string) string {
	if len(suggestions) == 0 {
		return ""
	}
	
	result := ""
	for i, suggestion := range suggestions {
		if i >= MaxSuggestionsDisplay {
			result += fmt.Sprintf("   ... and %d more", len(suggestions)-MaxSuggestionsDisplay)
			break
		}
		result += fmt.Sprintf("   • %s\n", suggestion)
	}
	return result[:len(result)-1] // Remove trailing newline
}