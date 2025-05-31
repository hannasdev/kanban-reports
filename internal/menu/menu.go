package menu

import (
	"bufio"
	"fmt"
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

// Menu handles interactive menu functionality
type Menu struct {
	scanner *bufio.Scanner
}

// NewMenu creates a new interactive menu
func NewMenu() *Menu {
	return &Menu{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// Run starts the interactive menu system
func (m *Menu) Run() (*config.Config, error) {
	fmt.Println("🔄 Kanban Reports - Interactive Mode")
	fmt.Println("=====================================")
	
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
	fmt.Println("\n📁 CSV File Selection")
	fmt.Println("--------------------")
	
	for {
		fmt.Print("Enter the path to your CSV file: ")
		if !m.scanner.Scan() {
			return "", fmt.Errorf("failed to read input")
		}
		
		path := strings.TrimSpace(m.scanner.Text())
		if path == "" {
			fmt.Println("❌ Please enter a valid file path")
			continue
		}
		
		// Perform comprehensive validation
		if err := validation.ValidateCSVPath(path); err != nil {
			csvErr, ok := err.(validation.CSVPathError)
			if !ok {
				fmt.Printf("❌ Error: %v\n", err)
				continue
			}
			
			// Handle different error types with helpful suggestions
			switch csvErr.Type {
			case "is_directory":
				fmt.Printf("❌ %s\n", csvErr.Message)
				
				// Suggest CSV files in the directory
				suggestions := validation.SuggestCSVFiles(path)
				if len(suggestions) > 0 {
					fmt.Println("\n💡 Found these CSV files in that directory:")
					for i, suggestion := range suggestions {
						if i >= 5 { // Limit suggestions
							fmt.Printf("   ... and %d more\n", len(suggestions)-5)
							break
						}
						fmt.Printf("   • %s\n", suggestion)
					}
					fmt.Println("\nPlease enter the full path to one of these files.")
				} else {
					fmt.Printf("\n💡 Try: %s/your-file.csv\n", path)
				}
				
			case "not_found":
				fmt.Printf("❌ %s\n", csvErr.Message)
				fmt.Println("💡 Make sure the file path is correct and the file exists.")
				
			case "not_readable":
				fmt.Printf("❌ %s\n", csvErr.Message)
				fmt.Println("💡 Check file permissions or if the file is open in another program.")
				
			case "empty_file":
				fmt.Printf("❌ %s\n", csvErr.Message)
				fmt.Println("💡 Make sure your CSV file contains data.")
				
			case "invalid_format":
				fmt.Printf("❌ %s\n", csvErr.Message)
				fmt.Println("💡 Make sure the file is a text-based CSV file, not binary.")
				
			default:
				fmt.Printf("❌ %s\n", csvErr.Message)
			}
			
			continue
		}
		
		fmt.Printf("✅ File validated: %s\n", path)
		return path, nil
	}
}

func (m *Menu) chooseMode() (bool, error) {
	fmt.Println("\n🎯 Mode Selection")
	fmt.Println("----------------")
	fmt.Println("Choose what you want to generate:")
	fmt.Println("1. 📊 Reports (story points by contributor, epic, team, or product area)")
	fmt.Println("2. 📈 Metrics (lead time, throughput, flow efficiency, etc.)")
	
	for {
		fmt.Print("\nEnter your choice (1 or 2): ")
		if !m.scanner.Scan() {
			return false, fmt.Errorf("failed to read input")
		}
		
		choice := strings.TrimSpace(m.scanner.Text())
		switch choice {
		case "1":
			return false, nil // Reports mode
		case "2":
			return true, nil // Metrics mode
		default:
			fmt.Println("❌ Please enter 1 or 2")
		}
	}
}

func (m *Menu) configureReports(cfg *config.Config) error {
	fmt.Println("\n📊 Report Type Selection")
	fmt.Println("------------------------")
	fmt.Println("Available report types:")
	fmt.Println("1. 👤 Contributor - Story points by person")
	fmt.Println("2. 🎯 Epic - Story points by epic/initiative")
	fmt.Println("3. 🏢 Product Area - Story points by product area")
	fmt.Println("4. 👥 Team - Story points by team")
	
	for {
		fmt.Print("\nEnter your choice (1-4): ")
		if !m.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}
		
		choice := strings.TrimSpace(m.scanner.Text())
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
			fmt.Println("❌ Please enter a number between 1 and 4")
			continue
		}
		
		cfg.ReportType = reportType
		fmt.Printf("✅ Selected: %s report\n", reportType)
		return nil
	}
}

func (m *Menu) configureMetrics(cfg *config.Config) error {
	fmt.Println("\n📈 Metrics Type Selection")
	fmt.Println("-------------------------")
	fmt.Println("Available metrics:")
	fmt.Println("1. ⏱️  Lead Time - How long items take to complete")
	fmt.Println("2. 🚀 Throughput - Completion rates over time")
	fmt.Println("3. 🌊 Flow Efficiency - Active vs waiting time")
	fmt.Println("4. 🎯 Estimation Accuracy - Estimate vs actual time correlation")
	fmt.Println("5. 📅 Work Item Age - Age of current incomplete items")
	fmt.Println("6. 📊 Team Improvement - Month-over-month trends")
	fmt.Println("7. 🔄 All Metrics - Generate all of the above")
	
	for {
		fmt.Print("\nEnter your choice (1-7): ")
		if !m.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}
		
		choice := strings.TrimSpace(m.scanner.Text())
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
		fmt.Printf("✅ Selected: %s metrics\n", metricsType)
		
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
	fmt.Println("\n⏰ Time Period Selection")
	fmt.Println("-----------------------")
	fmt.Println("Choose time period for grouping:")
	fmt.Println("1. 📅 Week - Group by week")
	fmt.Println("2. 🗓️  Month - Group by month")
	
	for {
		fmt.Print("\nEnter your choice (1 or 2): ")
		if !m.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}
		
		choice := strings.TrimSpace(m.scanner.Text())
		switch choice {
		case "1":
			cfg.PeriodType = metrics.PeriodTypeWeek
			fmt.Println("✅ Selected: Weekly grouping")
			return nil
		case "2":
			cfg.PeriodType = metrics.PeriodTypeMonth
			fmt.Println("✅ Selected: Monthly grouping")
			return nil
		default:
			fmt.Println("❌ Please enter 1 or 2")
		}
	}
}

func (m *Menu) configureDateRange(cfg *config.Config) error {
	fmt.Println("\n📅 Date Range Selection")
	fmt.Println("----------------------")
	fmt.Println("Choose date range:")
	fmt.Println("1. 🔄 All time - Include all data")
	fmt.Println("2. 📊 Last N days - Recent data only")
	fmt.Println("3. 📆 Specific range - Custom start and end dates")
	
	for {
		fmt.Print("\nEnter your choice (1-3): ")
		if !m.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}
		
		choice := strings.TrimSpace(m.scanner.Text())
		switch choice {
		case "1":
			fmt.Println("✅ Selected: All time")
			return nil
		case "2":
			return m.configureLastNDays(cfg)
		case "3":
			return m.configureSpecificRange(cfg)
		default:
			fmt.Println("❌ Please enter a number between 1 and 3")
		}
	}
}

func (m *Menu) configureLastNDays(cfg *config.Config) error {
	fmt.Println("\nCommon timeframes:")
	fmt.Println("- Last 7 days (1 week)")
	fmt.Println("- Last 30 days (1 month)")
	fmt.Println("- Last 90 days (1 quarter)")
	
	for {
		fmt.Print("\nEnter number of days: ")
		if !m.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}
		
		input := strings.TrimSpace(m.scanner.Text())
		days, err := strconv.Atoi(input)
		if err != nil || days <= 0 {
			fmt.Println("❌ Please enter a valid positive number")
			continue
		}
		
		cfg.LastNDays = days
		cfg.EndDate = time.Now()
		cfg.StartDate = cfg.EndDate.AddDate(0, 0, -days)
		
		fmt.Printf("✅ Selected: Last %d days\n", days)
		return nil
	}
}

func (m *Menu) configureSpecificRange(cfg *config.Config) error {
	// Get start date
	for {
		fmt.Print("\nEnter start date (YYYY-MM-DD): ")
		if !m.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}
		
		input := strings.TrimSpace(m.scanner.Text())
		startDate, err := time.Parse("2006-01-02", input)
		if err != nil {
			fmt.Println("❌ Invalid date format. Please use YYYY-MM-DD")
			continue
		}
		
		cfg.StartDate = startDate
		break
	}
	
	// Get end date
	for {
		fmt.Print("Enter end date (YYYY-MM-DD): ")
		if !m.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}
		
		input := strings.TrimSpace(m.scanner.Text())
		endDate, err := time.Parse("2006-01-02", input)
		if err != nil {
			fmt.Println("❌ Invalid date format. Please use YYYY-MM-DD")
			continue
		}
		
		// Add end of day to end date
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		
		if endDate.Before(cfg.StartDate) {
			fmt.Println("❌ End date cannot be before start date")
			continue
		}
		
		cfg.EndDate = endDate
		break
	}
	
	fmt.Printf("✅ Selected: %s to %s\n", 
		cfg.StartDate.Format("2006-01-02"), 
		cfg.EndDate.Format("2006-01-02"))
	return nil
}

func (m *Menu) configureFilters(cfg *config.Config) error {
	fmt.Println("\n🔍 Ad-hoc Request Filtering")
	fmt.Println("--------------------------")
	fmt.Println("How should ad-hoc requests be handled?")
	fmt.Println("1. ✅ Include all items (default)")
	fmt.Println("2. ❌ Exclude ad-hoc requests")
	fmt.Println("3. 🎯 Only ad-hoc requests")
	
	for {
		fmt.Print("\nEnter your choice (1-3): ")
		if !m.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}
		
		choice := strings.TrimSpace(m.scanner.Text())
		switch choice {
		case "1", "":
			cfg.AdHocFilter = "include"
			fmt.Println("✅ Selected: Include all items")
		case "2":
			cfg.AdHocFilter = "exclude"
			fmt.Println("✅ Selected: Exclude ad-hoc requests")
		case "3":
			cfg.AdHocFilter = "only"
			fmt.Println("✅ Selected: Only ad-hoc requests")
		default:
			fmt.Println("❌ Please enter a number between 1 and 3")
			continue
		}
		
		// Configure filter field
		cfg.FilterField = models.FilterFieldCompletedAt // Default
		return nil
	}
}

func (m *Menu) configureOutput(cfg *config.Config) error {
	fmt.Println("\n💾 Output Configuration")
	fmt.Println("----------------------")
	fmt.Println("Where should the report be displayed?")
	fmt.Println("1. 🖥️  Console only (display on screen)")
	fmt.Println("2. 📄 Save to file")
	
	for {
		fmt.Print("\nEnter your choice (1 or 2): ")
		if !m.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}
		
		choice := strings.TrimSpace(m.scanner.Text())
		switch choice {
		case "1":
			fmt.Println("✅ Selected: Console output")
			return nil
		case "2":
			return m.configureOutputFile(cfg)
		default:
			fmt.Println("❌ Please enter 1 or 2")
		}
	}
}

func (m *Menu) configureOutputFile(cfg *config.Config) error {
	for {
		fmt.Print("\nEnter output filename (e.g., report.txt): ")
		if !m.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}
		
		filename := strings.TrimSpace(m.scanner.Text())
		if filename == "" {
			fmt.Println("❌ Please enter a valid filename")
			continue
		}
		
		cfg.OutputPath = filename
		fmt.Printf("✅ Selected: Save to %s\n", filename)
		return nil
	}
}

func (m *Menu) configureDelimiter(cfg *config.Config) error {
	fmt.Println("\n🔗 CSV Delimiter Configuration")
	fmt.Println("-----------------------------")
	fmt.Println("Choose CSV delimiter (auto-detection recommended):")
	fmt.Println("1. 🤖 Auto-detect (recommended)")
	fmt.Println("2. , Comma")
	fmt.Println("3. ; Semicolon")
	fmt.Println("4. ⭾ Tab")
	
	for {
		fmt.Print("\nEnter your choice (1-4): ")
		if !m.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}
		
		choice := strings.TrimSpace(m.scanner.Text())
		switch choice {
		case "1", "":
			cfg.Delimiter = models.DelimiterAuto
			fmt.Println("✅ Selected: Auto-detection")
		case "2":
			cfg.Delimiter = models.DelimiterComma
			fmt.Println("✅ Selected: Comma delimiter")
		case "3":
			cfg.Delimiter = models.DelimiterSemicolon
			fmt.Println("✅ Selected: Semicolon delimiter")
		case "4":
			cfg.Delimiter = models.DelimiterTab
			fmt.Println("✅ Selected: Tab delimiter")
		default:
			fmt.Println("❌ Please enter a number between 1 and 4")
			continue
		}
		return nil
	}
}

// ShowSummary displays a summary of the selected configuration
func (m *Menu) ShowSummary(cfg *config.Config) {
	fmt.Println("\n📋 Configuration Summary")
	fmt.Println("=======================")
	fmt.Printf("📁 CSV File: %s\n", cfg.CSVPath)
	
	if cfg.IsMetricsReport() {
		fmt.Printf("📈 Metrics Type: %s\n", cfg.MetricsType)
		if cfg.MetricsType == metrics.MetricsTypeThroughput || cfg.MetricsType == metrics.MetricsTypeAll {
			fmt.Printf("⏰ Period: %s\n", cfg.PeriodType)
		}
	} else {
		fmt.Printf("📊 Report Type: %s\n", cfg.ReportType)
	}
	
	// Date range
	if cfg.LastNDays > 0 {
		fmt.Printf("📅 Date Range: Last %d days\n", cfg.LastNDays)
	} else if !cfg.StartDate.IsZero() && !cfg.EndDate.IsZero() {
		fmt.Printf("📅 Date Range: %s to %s\n", 
			cfg.StartDate.Format("2006-01-02"), 
			cfg.EndDate.Format("2006-01-02"))
	} else {
		fmt.Printf("📅 Date Range: All time\n")
	}
	
	fmt.Printf("🔍 Ad-hoc Filter: %s\n", cfg.AdHocFilter)
	fmt.Printf("🔗 Delimiter: %s\n", cfg.Delimiter.Name)
	
	if cfg.OutputPath != "" {
		fmt.Printf("💾 Output: %s\n", cfg.OutputPath)
	} else {
		fmt.Printf("💾 Output: Console\n")
	}
	
	fmt.Println("\n🚀 Generating report...")
}