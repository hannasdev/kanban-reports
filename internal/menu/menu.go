package menu

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hannasdev/kanban-reports/internal/config"
	"github.com/hannasdev/kanban-reports/internal/metrics"
	"github.com/hannasdev/kanban-reports/internal/models"
	"github.com/hannasdev/kanban-reports/internal/reports"
	"github.com/hannasdev/kanban-reports/internal/validation"
)

// MenuInterface defines the interface for input/output operations
type MenuInterface interface {
	ReadInput(prompt string) (string, error)
	Print(msg string)
	Printf(format string, args ...interface{})
	Println(msg string)
}

// Menu handles interactive menu functionality
type Menu struct {
	scanner *bufio.Scanner
	writer  io.Writer
	reader  io.Reader
}

// NewMenu creates a new interactive menu
func NewMenu() *Menu {
	return &Menu{
		scanner: bufio.NewScanner(os.Stdin),
		writer:  os.Stdout,
		reader:  os.Stdin,
	}
}

// NewMenuWithIO creates a new menu with custom input/output for testing
func NewMenuWithIO(reader io.Reader, writer io.Writer) *Menu {
	return &Menu{
		scanner: bufio.NewScanner(reader),
		writer:  writer,
		reader:  reader,
	}
}

func (m *Menu) print(msg string) {
	fmt.Fprint(m.writer, msg)
}

// printf outputs formatted text to the configured writer
func (m *Menu) printf(format string, args ...interface{}) {
	fmt.Fprintf(m.writer, format, args...)
}

// println outputs a line to the configured writer
func (m *Menu) println(msg string) {
	fmt.Fprintln(m.writer, msg)
}

// readInput reads input from user and checks for quit commands
func (m *Menu) readInput(prompt string) (string, error) {
	m.print(prompt)
	if !m.scanner.Scan() {
		return "", fmt.Errorf("failed to read input")
	}
	
	input := m.scanner.Text()
	
	// Check for quit command
	if err := HandleQuit(input); err != nil {
		return "", err
	}
	
	return strings.TrimSpace(input), nil
}

// Run starts the interactive menu system
func (m *Menu) Run() (*config.Config, error) {
	m.println("🔄 Kanban Reports - Interactive Mode")
	m.println("=====================================")
	ShowQuitHelp()
	
	cfg := &config.Config{}
	
	// Step 1: Get CSV file path
	csvPath, err := m.getCSVPath()
	if err != nil {
		return nil, err
	}
	cfg.CSVPath = csvPath
	
	// Step 2: Choose report or metrics mode
	isMetrics, err := m.chooseMode()
	if err != nil {
		return nil, err
	}
	
	if isMetrics {
		// Step 3a: Configure metrics
		if err := m.configureMetrics(cfg); err != nil {
			return nil, err
		}
	} else {
		// Step 3b: Configure reports
		if err := m.configureReports(cfg); err != nil {
			return nil, err
		}
	}
	
	// Step 4: Configure date range
	if err := m.configureDateRange(cfg); err != nil {
		return nil, err
	}
	
	// Step 5: Configure filters
	if err := m.configureFilters(cfg); err != nil {
		return nil, err
	}
	
	// Step 6: Configure output
	if err := m.configureOutput(cfg); err != nil {
		return nil, err
	}
	
	// Step 7: Configure delimiter
	if err := m.configureDelimiter(cfg); err != nil {
		return nil, err
	}
	
	return cfg, nil
}

func (m *Menu) getCSVPath() (string, error) {
	m.println("\n📁 CSV File Selection")
	m.println("--------------------")
	
	for {
		path, err := m.readInput("Enter the path to your CSV file: ")
		if err != nil {
			// This already handles quit commands from readInput
			return "", err
		}
		
		if path == "" {
			m.println("❌ Please enter a valid file path")
			continue
		}
		
		// Perform comprehensive validation
		if err := validation.ValidateCSVPath(path); err != nil {
			csvErr, ok := err.(validation.CSVPathError)
			if !ok {
				m.printf("❌ Error: %v\n", err)
				continue
			}
			
			// Handle different error types with helpful suggestions
			switch csvErr.Type {
			case "is_directory":
				m.printf("❌ %s\n", csvErr.Message)
				
				// Suggest CSV files in the directory
				suggestions := validation.SuggestCSVFiles(path)
				if len(suggestions) > 0 {
					m.println("\n💡 Found these CSV files in that directory:")
					for i, suggestion := range suggestions {
						if i >= 5 { // Limit suggestions
							m.printf("   • %s\n", suggestion)
						} else {
							m.printf("   ... and %d more\n", len(suggestions)-5)
							break
						}
					}
					m.println("\nPlease enter the full path to one of these files.")
				} else {
					m.printf("\n💡 Try: %s/your-file.csv\n", path)
				}
				
			case "not_found":
				m.printf("❌ %s\n", csvErr.Message)
				m.println("💡 Make sure the file path is correct and the file exists.")
				
			case "not_readable":
				m.printf("❌ %s\n", csvErr.Message)
				m.println("💡 Check file permissions or if the file is open in another program.")
				
			case "empty_file":
				m.printf("❌ %s\n", csvErr.Message)
				m.println("💡 Make sure your CSV file contains data.")
				
			case "invalid_format":
				m.printf("❌ %s\n", csvErr.Message)
				m.println("💡 Make sure the file is a text-based CSV file, not binary.")
				
			default:
				m.printf("❌ %s\n", csvErr.Message)
			}
			
			// The continue here will loop back to readInput, which handles quit
			continue
		}
		
		m.printf("✅ File validated: %s\n", path)
		return path, nil
	}
}

func (m *Menu) chooseMode() (bool, error) {
	m.println("\n🎯 Mode Selection")
	m.println("----------------")
	m.println("Choose what you want to generate:")
	m.println("1. 📊 Reports (story points by contributor, epic, team, or product area)")
	m.println("2. 📈 Metrics (lead time, throughput, flow efficiency, etc.)")
	
	for {
		choice, err := m.readInput("\nEnter your choice (1 or 2): ")
		if err != nil {
			return false, err
		}
		
		switch choice {
		case "1":
			return false, nil // Reports mode
		case "2":
			return true, nil // Metrics mode
		default:
			m.println("❌ Please enter 1 or 2")
		}
	}
}

func (m *Menu) configureReports(cfg *config.Config) error {
	m.println("\n📊 Report Type Selection")
	m.println("------------------------")
	m.println("Available report types:")
	m.println("1. 👤 Contributor - Story points by person")
	m.println("2. 🎯 Epic - Story points by epic/initiative")
	m.println("3. 🏢 Product Area - Story points by product area")
	m.println("4. 👥 Team - Story points by team")
	
	for {
		choice, err := m.readInput("\nEnter your choice (1-4): ")
		if err != nil {
			return err
		}
		
		var reportType reports.ReportType
		
		switch choice {
		case "1":
			reportType = reports.ReportTypeContributor
		case "2":
			reportType = reports.ReportTypeEpic
		case "3":
			reportType = reports.ReportTypeProductArea
		case "4":
			reportType = reports.ReportTypeTeam
		default:
			m.println("❌ Please enter a number between 1 and 4")
			continue
		}
		
		cfg.ReportType = reportType
		m.printf("✅ Selected: %s report\n", reportType)
		return nil
	}
}

func (m *Menu) configureMetrics(cfg *config.Config) error {
	m.println("\n📈 Metrics Type Selection")
	m.println("-------------------------")
	m.println("Available metrics:")
	m.println("1. ⏱️  Lead Time - How long items take to complete")
	m.println("2. 🚀 Throughput - Completion rates over time")
	m.println("3. 🌊 Flow Efficiency - Active vs waiting time")
	m.println("4. 🎯 Estimation Accuracy - Estimate vs actual time correlation")
	m.println("5. 📅 Work Item Age - Age of current incomplete items")
	m.println("6. 📊 Team Improvement - Month-over-month trends")
	m.println("7. 🔄 All Metrics - Generate all of the above")
	
	for {
		choice, err := m.readInput("\nEnter your choice (1-7): ")
		if err != nil {
			return err
		}
		
		var metricsType metrics.MetricsType
		
		switch choice {
		case "1":
			metricsType = metrics.MetricsTypeLeadTime
		case "2":
			metricsType = metrics.MetricsTypeThroughput
		case "3":
			metricsType = metrics.MetricsTypeFlow
		case "4":
			metricsType = metrics.MetricsTypeEstimation
		case "5":
			metricsType = metrics.MetricsTypeAge
		case "6":
			metricsType = metrics.MetricsTypeImprovement
		case "7":
			metricsType = metrics.MetricsTypeAll
		default:
			fmt.Println("❌ Please enter a number between 1 and 7")
			continue
		}
		
		cfg.MetricsType = metricsType
		m.printf("✅ Selected: %s metrics\n", metricsType)
		
		// For throughput metrics, ask about period
		if metricsType == metrics.MetricsTypeThroughput || metricsType == metrics.MetricsTypeAll {
			return m.configurePeriod(cfg)
		}
		
		// Set default period for other metrics
		cfg.PeriodType = metrics.PeriodTypeMonth
		return nil
	}
}

func (m *Menu) configurePeriod(cfg *config.Config) error {
	m.println("\n⏰ Time Period Selection")
	m.println("-----------------------")
	m.println("Choose time period for grouping:")
	m.println("1. 📅 Week - Group by week")
	m.println("2. 🗓️  Month - Group by month")
	
	for {
		choice, err := m.readInput("\nEnter your choice (1 or 2): ")
		if err != nil {
			return err
		}
		
		switch choice {
		case "1":
			cfg.PeriodType = metrics.PeriodTypeWeek
			m.println("✅ Selected: Weekly grouping")
			return nil
		case "2":
			cfg.PeriodType = metrics.PeriodTypeMonth
			m.println("✅ Selected: Monthly grouping")
			return nil
		default:
			m.println("❌ Please enter 1 or 2")
		}
	}
}

func (m *Menu) configureDateRange(cfg *config.Config) error {
	m.println("\n📅 Date Range Selection")
	m.println("----------------------")
	m.println("Choose date range:")
	m.println("1. 🔄 All time - Include all data")
	m.println("2. 📊 Last N days - Recent data only")
	m.println("3. 📆 Specific range - Custom start and end dates")
	
	for {
		choice, err := m.readInput("\nEnter your choice (1-3): ")
		if err != nil {
			return err
		}
		
		switch choice {
		case "1":
			m.println("✅ Selected: All time")
			return nil
		case "2":
			return m.configureLastNDays(cfg)
		case "3":
			return m.configureSpecificRange(cfg)
		default:
			m.println("❌ Please enter a number between 1 and 3")
		}
	}
}

func (m *Menu) configureLastNDays(cfg *config.Config) error {
	m.println("\nCommon timeframes:")
	m.println("- Last 7 days (1 week)")
	m.println("- Last 30 days (1 month)")
	m.println("- Last 90 days (1 quarter)")
	
	for {
		input, err := m.readInput("\nEnter number of days: ")
		if err != nil {
			return err
		}
		
		days, err := strconv.Atoi(input)
		if err != nil || days <= 0 {
			m.println("❌ Please enter a valid positive number")
			continue
		}
		
		cfg.LastNDays = days
		cfg.EndDate = time.Now()
		cfg.StartDate = cfg.EndDate.AddDate(0, 0, -days)
		
		m.printf("✅ Selected: Last %d days\n", days)
		return nil
	}
}

func (m *Menu) configureSpecificRange(cfg *config.Config) error {
	// Get start date
	for {
		input, err := m.readInput("\nEnter start date (YYYY-MM-DD): ")
		if err != nil {
			return err
		}
		
		startDate, err := time.Parse("2006-01-02", input)
		if err != nil {
			m.println("❌ Invalid date format. Please use YYYY-MM-DD")
			continue
		}
		
		cfg.StartDate = startDate
		break
	}
	
	// Get end date
	for {
		input, err := m.readInput("Enter end date (YYYY-MM-DD): ")
		if err != nil {
			return err
		}
		
		endDate, err := time.Parse("2006-01-02", input)
		if err != nil {
			m.println("❌ Invalid date format. Please use YYYY-MM-DD")
			continue
		}
		
		// Add end of day to end date
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		
		if endDate.Before(cfg.StartDate) {
			m.println("❌ End date cannot be before start date")
			continue
		}
		
		cfg.EndDate = endDate
		break
	}
	
	m.printf("✅ Selected: %s to %s\n", 
		cfg.StartDate.Format("2006-01-02"), 
		cfg.EndDate.Format("2006-01-02"))
	return nil
}

func (m *Menu) configureFilters(cfg *config.Config) error {
	m.println("\n🔍 Ad-hoc Request Filtering")
	m.println("--------------------------")
	m.println("How should ad-hoc requests be handled?")
	m.println("1. ✅ Include all items (default)")
	m.println("2. ❌ Exclude ad-hoc requests")
	m.println("3. 🎯 Only ad-hoc requests")
	
	for {
		choice, err := m.readInput("\nEnter your choice (1-3): ")
		if err != nil {
			return err
		}
		
		switch choice {
		case "1", "":
			cfg.AdHocFilter = "include"
			m.println("✅ Selected: Include all items")
		case "2":
			cfg.AdHocFilter = "exclude"
			m.println("✅ Selected: Exclude ad-hoc requests")
		case "3":
			cfg.AdHocFilter = "only"
			m.println("✅ Selected: Only ad-hoc requests")
		default:
			m.println("❌ Please enter a number between 1 and 3")
			continue
		}
		
		// Configure filter field
		cfg.FilterField = models.FilterFieldCompletedAt // Default
		return nil
	}
}

func (m *Menu) configureOutput(cfg *config.Config) error {
	m.println("\n💾 Output Configuration")
	m.println("----------------------")
	m.println("Where should the report be displayed?")
	m.println("1. 🖥️  Console only (display on screen)")
	m.println("2. 📄 Save to file")
	
	for {
		choice, err := m.readInput("\nEnter your choice (1 or 2): ")
		if err != nil {
			return err
		}
		
		switch choice {
		case "1":
			m.println("✅ Selected: Console output")
			return nil
		case "2":
			return m.configureOutputFile(cfg)
		default:
			m.println("❌ Please enter 1 or 2")
		}
	}
}

func (m *Menu) configureOutputFile(cfg *config.Config) error {
	for {
		filename, err := m.readInput("\nEnter output filename (e.g., report.txt): ")
		if err != nil {
			return err
		}
		
		if filename == "" {
			m.println("❌ Please enter a valid filename")
			continue
		}
		
		cfg.OutputPath = filename
		m.printf("✅ Selected: Save to %s\n", filename)
		return nil
	}
}

func (m *Menu) configureDelimiter(cfg *config.Config) error {
	m.println("\n🔗 CSV Delimiter Configuration")
	m.println("-----------------------------")
	m.println("Choose CSV delimiter (auto-detection recommended):")
	m.println("1. 🤖 Auto-detect (recommended)")
	m.println("2. , Comma")
	m.println("3. ; Semicolon")
	m.println("4. ⭾ Tab")
	
	for {
		choice, err := m.readInput("\nEnter your choice (1-4): ")
		if err != nil {
			return err
		}
		
		switch choice {
		case "1", "":
			cfg.Delimiter = models.DelimiterAuto
			m.println("✅ Selected: Auto-detection")
		case "2":
			cfg.Delimiter = models.DelimiterComma
			m.println("✅ Selected: Comma delimiter")
		case "3":
			cfg.Delimiter = models.DelimiterSemicolon
			m.println("✅ Selected: Semicolon delimiter")
		case "4":
			cfg.Delimiter = models.DelimiterTab
			m.println("✅ Selected: Tab delimiter")
		default:
			m.println("❌ Please enter a number between 1 and 4")
			continue
		}
		return nil
	}
}

// ShowSummary displays a summary of the selected configuration
func (m *Menu) ShowSummary(cfg *config.Config) {
	m.println("\n📋 Configuration Summary")
	m.println("=======================")
	m.printf("📁 CSV File: %s\n", cfg.CSVPath)
	
	if cfg.IsMetricsReport() {
		m.printf("📈 Metrics Type: %s\n", cfg.MetricsType)
		if cfg.MetricsType == metrics.MetricsTypeThroughput || cfg.MetricsType == metrics.MetricsTypeAll {
			m.printf("⏰ Period: %s\n", cfg.PeriodType)
		}
	} else {
		m.printf("📊 Report Type: %s\n", cfg.ReportType)
	}
	
	// Date range
	if cfg.LastNDays > 0 {
		m.printf("📅 Date Range: Last %d days\n", cfg.LastNDays)
	} else if !cfg.StartDate.IsZero() && !cfg.EndDate.IsZero() {
		m.printf("📅 Date Range: %s to %s\n", 
			cfg.StartDate.Format("2006-01-02"), 
			cfg.EndDate.Format("2006-01-02"))
	} else {
		m.printf("📅 Date Range: All time\n")
	}
	
	m.printf("🔍 Ad-hoc Filter: %s\n", cfg.AdHocFilter)
	m.printf("🔗 Delimiter: %s\n", cfg.Delimiter.Name)
	
	if cfg.OutputPath != "" {
		m.printf("💾 Output: %s\n", cfg.OutputPath)
	} else {
		m.printf("💾 Output: Console\n")
	}
	
	m.println("\n🚀 Generating report...")
}