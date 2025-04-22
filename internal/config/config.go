package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/hannasdev/kanban-reports/internal/metrics"
	"github.com/hannasdev/kanban-reports/internal/reports"
)

// Config represents the application configuration
type Config struct {
	// Input file configuration
	CSVPath     string
	Delimiter   rune
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
	AdHocFilter reports.AdHocFilterType
}

// ParseFlags parses command-line flags and returns a populated Config
func ParseFlags() (*Config, error) {
	config := &Config{}

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

	// Validate and set CSV path
	if *csvPath == "" {
		return nil, fmt.Errorf("CSV file path is required")
	}
	config.CSVPath = *csvPath

	// Set delimiter
	config.AutoDetect = true
	switch *delimiterStr {
	case "comma":
		config.Delimiter = ','
		config.AutoDetect = false
	case "tab":
		config.Delimiter = '\t'
		config.AutoDetect = false
	case "semicolon":
		config.Delimiter = ';'
		config.AutoDetect = false
	case "auto":
		// Auto-detection is the default
	default:
		fmt.Println("Invalid delimiter specified, using auto-detection")
	}

	// Parse report type
	if *metricsType == "" {
		switch *reportType {
		case "contributor":
			config.ReportType = reports.ReportTypeContributor
		case "epic":
			config.ReportType = reports.ReportTypeEpic
		case "product-area":
			config.ReportType = reports.ReportTypeProductArea
		case "team":
			config.ReportType = reports.ReportTypeTeam
		default:
			return nil, fmt.Errorf("unknown report type: %s", *reportType)
		}
	}

	// Parse metrics type
	switch *metricsType {
	case "":
		// No metrics specified, using report type
	case "lead-time":
		config.MetricsType = metrics.MetricsTypeLeadTime
	case "throughput":
		config.MetricsType = metrics.MetricsTypeThroughput
	case "flow":
		config.MetricsType = metrics.MetricsTypeFlow
	case "estimation":
		config.MetricsType = metrics.MetricsTypeEstimation
	case "age":
		config.MetricsType = metrics.MetricsTypeAge
	case "improvement":
		config.MetricsType = metrics.MetricsTypeImprovement
	case "all":
		config.MetricsType = metrics.MetricsTypeAll
	default:
		return nil, fmt.Errorf("unknown metrics type: %s", *metricsType)
	}

	// Parse period type
	switch *periodType {
	case "week":
		config.PeriodType = metrics.PeriodTypeWeek
	case "month":
		config.PeriodType = metrics.PeriodTypeMonth
	default:
		return nil, fmt.Errorf("unknown period type: %s", *periodType)
	}

	// Parse date range
	var err error

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

	// Parse ad-hoc filter
	switch *adHocFilter {
	case "include":
		config.AdHocFilter = reports.AdHocFilterInclude
	case "exclude":
		config.AdHocFilter = reports.AdHocFilterExclude
	case "only":
		config.AdHocFilter = reports.AdHocFilterOnly
	default:
		return nil, fmt.Errorf("unknown ad-hoc filter type: %s", *adHocFilter)
	}

	// Set output path
	config.OutputPath = *outputPath

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