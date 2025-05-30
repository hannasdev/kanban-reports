package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/hannasdev/kanban-reports/internal/metrics"
	"github.com/hannasdev/kanban-reports/internal/models"
	"github.com/hannasdev/kanban-reports/internal/reports"
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
}

// ParseFlags parses command-line flags and returns a populated Config
func ParseFlags() (*Config, error) {
	config := &Config{}
	var err error

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
	filterField := flag.String("filter-field", "completed_at", "Date field to filter by: completed_at, created_at, started_at")


	flag.Parse()

	// Validate and set CSV path
	if *csvPath == "" {
		return nil, fmt.Errorf("CSV file path is required")
	}
	config.CSVPath = *csvPath

	// Set delimiter
	config.Delimiter, err = models.ParseDelimiter(*delimiterStr)
	if err != nil {
			fmt.Printf("Warning: %v, using auto-detection\n", err)
			config.Delimiter = models.DelimiterAuto
	}

	// Parse report type with validation
	if *metricsType == "" {
		config.ReportType, err = reports.ParseReportType(*reportType)
		if err != nil {
				return nil, err
		}
	}

	// Parse metrics type with validation
	if *metricsType != "" {
		config.MetricsType, err = metrics.ParseMetricsType(*metricsType)
		if err != nil {
				return nil, err
		}
	}

	// Parse period type with validation
	config.PeriodType, err = metrics.ParsePeriodType(*periodType)
	if err != nil {
			return nil, err
	}

	// Parse ad-hoc filter with validation
	config.AdHocFilter, err = types.ParseAdHocFilterType(*adHocFilter)
	if err != nil {
			return nil, err
	}

	 // Validate filter field
	 config.FilterField, err = models.ParseFilterField(*filterField)
	 if err != nil {
			 return nil, err
	 }

	 // Validate and process date range
	if *lastNDays < 0 {
		return nil, fmt.Errorf("last N days must be a positive number, got: %d", *lastNDays)
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
						return nil, fmt.Errorf("error parsing start date: %v", err)
				}
		}

		if *endDateStr != "" {
				config.EndDate, err = time.Parse("2006-01-02", *endDateStr)
				if err != nil {
						return nil, fmt.Errorf("error parsing end date: %v", err)
				}
				// Set end date to the end of the day
				config.EndDate = config.EndDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
	}

	// Set output path
	config.OutputPath = *outputPath

	// Validate date range consistency
	if !config.StartDate.IsZero() && !config.EndDate.IsZero() && config.EndDate.Before(config.StartDate) {
		return nil, fmt.Errorf("invalid date range: end date (%s) is before start date (%s)", 
						config.EndDate.Format("2006-01-02"), config.StartDate.Format("2006-01-02"))
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