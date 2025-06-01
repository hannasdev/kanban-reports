package config

import (
	"fmt"
	"os"
)

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
		// Provide all 24 arguments for the format placeholders
		os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], 
		os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], 
		os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], 
		os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

// getGoVersion returns the Go version for version display
func getGoVersion() string {
	// In a real implementation, you might want to embed this at build time
	return "1.21+"
}

// getPlatform returns the platform information for version display
func getPlatform() string {
	// In a real implementation, you might want to embed this at build time
	return "linux/amd64"
}