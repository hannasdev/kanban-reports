# Kanban Reports

A Go application for generating reports from Kanban board data exported as CSV. This tool helps track and analyze team productivity by generating reports on story points delivered, broken down by contributor, epic, product area, or team.

## Features

- Process CSV exports from kanban boards
- Generate reports by:
  - Contributor (who completed the work)
  - Epic (which larger initiatives the work belongs to)
  - Product Area (which product categories the work affects)
  - Team (which teams are delivering the most)
- Generate advanced metrics:
  - Lead Time Analysis (how long items take to complete)
  - Throughput Analysis (completion rates over time)
  - Flow Efficiency (active vs. waiting time)
  - Estimation Accuracy (correlation between estimates and actual time)
  - Work Item Age (analysis of current incomplete work)
  - Team Improvement (month-over-month trends)
- Filter by date ranges or last N days
- Filter ad-hoc requests (include, exclude, or focus only on them)
- Save reports to file or view in console
- Automatic CSV delimiter detection (comma, tab, semicolon)

## Quick Start

### Option 1: Interactive Mode (Recommended for beginners)

#### Build and run setup

./scripts/setup.sh

#### Start interactive mode

./bin/kanban-reports --interactive

### Option 2: Command Line Mode

#### Show all available options

./bin/kanban-reports --help

#### See practical examples

./bin/kanban-reports --examples

#### Quick start with sample data

./bin/kanban-reports --csv data/sample.csv --type contributor --last 30

## Installation

### Prerequisites

- Go 1.21 or higher

### From Source

```bash
# Clone the repository
git clone https://github.com/hannasdev/kanban-reports.git
cd kanban-reports

# Build the application
go build -o bin/kanban-reports ./cmd/kanban-reports
```

## Usage

### Basic Usage

```bash
# Generate a contributor report for the last 7 days
./bin/kanban-reports --csv data/kanban-data.csv --type contributor --last 7

# Generate an epic report for a specific date range
./bin/kanban-reports --csv data/kanban-data.csv --type epic --start 2024-05-01 --end 2024-05-31

# Save the report to a file
./bin/kanban-reports --csv data/kanban-data.csv --type product-area --last 14 --output product-report.txt

# Generate lead time metrics by story point size
./bin/kanban-reports --csv data/kanban-data.csv --metrics lead-time --last 90

# Generate throughput metrics by week
./bin/kanban-reports --csv data/kanban-data.csv --metrics throughput --period week --last 180

# Generate all metrics for a specific date range
./bin/kanban-reports --csv data/kanban-data.csv --metrics all --start 2024-01-01 --end 2024-06-30 --output all-metrics.txt

# Generate reports excluding ad-hoc requests
./bin/kanban-reports --csv data/kanban-data.csv --type contributor --last 30 --ad-hoc exclude

# Analyze only ad-hoc requests
./bin/kanban-reports --csv data/kanban-data.csv --metrics throughput --last 90 --ad-hoc only
```

### Command Line Options

| Flag | Description | Example |
|------|-------------|---------|
| `--help, -h` | Complete help with all options | `./bin/kanban-reports --help` |
| `--examples` | Practical usage examples | `./bin/kanban-reports --examples` |
| `--version` | Version information | `./bin/kanban-reports --version` |
| `--interactive, -i` | Interactive menu mode | `./bin/kanban-reports` |
| `--csv` | Path to the kanban CSV file (required) | `--csv data/kanban-data.csv` |
| `--type` | Type of report to generate (contributor, epic, product-area, team) | `--type epic` |
| `--metrics` | Type of metrics to generate (lead-time, throughput, flow, estimation, age, improvement, all) | `--metrics lead-time` |
| `--period` | Time period for metrics reports (week, month) | `--period week` |
| `--start` | Start date in YYYY-MM-DD format | `--start 2024-05-01` |
| `--end` | End date in YYYY-MM-DD format | `--end 2024-05-31` |
| `--last` | Generate report for the last N days | `--last 7` |
| `--output` | Path to save the report (optional) | `--output report.txt` |
| `--delimiter` | CSV delimiter: comma, tab, semicolon, or auto (default: auto) | `--delimiter comma` |
| `--ad-hoc` | How to handle ad-hoc requests: include, exclude, only (default: include) | `--ad-hoc exclude` |

## CSV Data Format

The application expects a CSV file with the following required columns:

- `id`: Unique identifier for the item
- `name`: Name or title of the item
- `estimate`: Story points or estimate value
- `is_completed`: Boolean indicating if the item is completed
- `completed_at`: Date when the item was completed

Additional useful columns:

- `owners`: Person(s) assigned to the item
- `epic`: Epic name
- `team`: Team name
- `product_area`: Product area or category
- `type`: Type of work (feature, bug, task, etc.)

Example of the first row of CSV data:

```csv
id,name,type,requester,owners,description,is_completed,created_at,started_at,updated_at,moved_at,completed_at,estimate,external_ticket_count,external_tickets,is_blocked,is_a_blocker,due_date,labels,epic_labels,tasks,state,epic_id,epic,project_id,project,iteration_id,iteration,utc_offset,is_archived,team_id,team,epic_state,epic_is_archived,epic_created_at,epic_started_at,epic_due_date,milestone_id,milestone,milestone_state,milestone_created_at,milestone_started_at,milestone_due_date,milestone_categories,epic_planned_start_date,workflow,workflow_id,priority,severity,product_area,skill_set,technical_area,custom_fields
```

## Report Types

### Epic Report

Shows story points completed by epic.

Example output:

```zsh
Story Points by Epic:

API Modernization                      20.0 points   7 items
User Authentication Overhaul           15.5 points   5 items
No Epic                                 8.0 points   3 items

Total: 43.5 points across 15 items
```

### Product Area Report

Shows story points completed by product area.

Example output:

```zsh
Story Points by Product Area:

Backend                         25.0 points  10 items
Frontend                        18.5 points   7 items
Uncategorized                    5.0 points   2 items

Total: 48.5 points across 19 items
```

### Team Report

Shows story points completed by team.

Example output:

```zsh
Story Points by Team:

Team Alpha                      30.0 points  12 items
Team Beta                       20.5 points   8 items
No Team                          4.0 points   2 items

Total: 54.5 points across 22 items
```

## Metrics Types

### Lead Time Metrics

Analyzes how long items take to complete, broken down by story point size.

Example output:

```zsh
# Lead Time Analysis by Story Point Size (in days)

## Lead Time (Creation to Completion)

Story points | Count | Min | Max | Avg | Median
-------------|-------|-----|-----|-----|-------
           1 |    15 | 2.5 | 8.3 | 4.2 |    3.8
           3 |    12 | 3.1 | 9.7 | 5.8 |    5.2
           5 |     8 | 5.7 | 15.2 | 10.1 |    9.5
           8 |     5 | 7.2 | 21.5 | 14.3 |   12.8
```

### Throughput Metrics

Shows completion rates over time, broken down by week or month.

Example output:

```zsh
# Throughput Analysis by Month

Month   | Items Completed | Story Points | Avg Points/Item
--------|----------------|-------------|---------------
2024-01 |             12 |        45.0 |           3.8
2024-02 |             15 |        62.0 |           4.1
2024-03 |             18 |        72.0 |           4.0
```

### Flow Efficiency Metrics

Analyzes the percentage of time items spend in active states versus waiting.

Example output:

```zsh
# Flow Efficiency Analysis

State   | Avg Time (days) | % of Total Time
--------|-----------------|---------------
Waiting |            12.5 |         71.4%
Active  |             5.0 |         28.6%

Flow Efficiency: 28.6%
```

### Estimation Accuracy Metrics

Compares story point estimates with actual completion times.

Example output:

```zsh
# Estimation Accuracy Analysis

## Time Spent per Story Point Size

Story points | Count | Min Days/SP | Max Days/SP | Avg Days/SP | Median Days/SP
-------------|-------|------------|------------|-------------|---------------
           1 |    15 |        1.2 |        3.5 |         2.1 |             1.9
           3 |    12 |        1.0 |        2.8 |         1.8 |             1.7
           5 |     8 |        1.1 |        2.5 |         1.7 |             1.6
```

### Work Item Age Metrics

Analyzes the age of current work items by state.

Example output:

```zsh
# Current Work Item Age Analysis

## In Progress (5 items)

Min: 2.1, Max: 15.3, Avg: 7.2, Median: 5.5 days

Oldest Items:

- API Authentication Refactoring (15.3 days)
- Payment Processing Bug (10.8 days)
- User Profile UI Update (7.4 days)
```

### Team Improvement Metrics

Shows trends in team performance over time.

Example output:

```zsh
# Team Improvement Metrics

Month   | Items | Points | Avg Lead Time | Avg Cycle Time | Lead Time Δ | Cycle Time Δ
--------|-------|--------|---------------|----------------|------------|-------------
2024-01 |    12 |   45.0 |          18.5 |            8.2 |            |             
2024-02 |    15 |   62.0 |          16.8 |            7.5 | -1.7 (-9.2%) | -0.7 (-8.5%)
2024-03 |    18 |   72.0 |          14.2 |            6.8 | -2.6 (-15.5%) | -0.7 (-9.3%)
```

## Project Structure

```zsh
kanban-reports/
├── cmd/
│   └── kanban-reports/         # Main application entry point
│       └── main.go
├── internal/
│   ├── config/                 # Application configuration
│   │   └── config.go
│   ├── models/                 # Data models
│   │   └── kanban_item.go
│   ├── parser/                 # CSV parsing logic
│   │   └── csv_parser.go
│   ├── reports/                # Report generation
│   │   ├── reporter.go         # Core reporter functionality
│   │   ├── types.go            # Type definitions
│   │   ├── contributor.go      # Contributor report implementation
│   │   ├── epic.go             # Epic report implementation
│   │   ├── product_area.go     # Product area report implementation
│   │   └── team.go             # Team report implementation
│   └── metrics/                # Metrics generation
│       ├── metrics.go          # Core metrics generator
│       ├── types.go            # Metrics type definitions
│       ├── lead_time.go        # Lead time metrics
│       ├── throughput.go       # Throughput metrics
│       ├── flow.go             # Flow efficiency metrics
│       ├── estimation.go       # Estimation accuracy metrics
│       ├── age.go              # Work item age metrics
│       ├── improvement.go      # Team improvement metrics
│       └── util.go             # Common utility functions for metrics
├── pkg/                        # Reusable public packages
│   └── dateutil/               # Date handling utilities
│       └── dateutil.go
├── test/                       # Test data and helpers
│   └── fixtures/               # Test fixture files
│       └── sample.csv
├── data/                      # Place your CSV files here
│   └── kanban-data.csv
├── go.mod
└── go.sum
```

## Development

### Testing the Application

Run the unit tests to ensure everything is working properly:

```bash
# Run all unit tests
go test ./internal/...

# Run tests for a specific package
go test ./internal/reports/...

# Run with verbose output
go test -v ./internal/...

# Run with test coverage report
go test -cover ./internal/...
```

### Adding New Report Types

To add a new report type:

1. Add a new report type constant in `internal/reports/types.go`
2. Create a new file in `internal/reports/` for your report implementation
3. Add a new report generation function in the Reporter struct
4. Update the switch statement in `GenerateReport` to handle the new type
5. Update `cmd/kanban-reports/main.go` to accept the new report type as a command-line option

### Adding New Metrics

To add a new metrics type:

1. Add a new metrics type constant in `internal/metrics/types.go`
2. Create a new file in `internal/metrics/` for your metrics implementation
3. Add a new metrics generation function
4. Update the switch statement in `Generator.Generate()` to handle the new type
5. Update `cmd/kanban-reports/main.go` to accept the new metrics type as a command-line option

### Modifying the CSV Parser

If your CSV format changes, you may need to update:

1. The `KanbanItem` struct in `internal/models/kanban_item.go`
2. The parsing logic in `internal/parser/csv_parser.go`

## Example Workflow

1. Export your kanban data as a CSV file from your kanban board tool
2. Place the CSV file in the `data/` directory
3. Run the application with your desired report type and date range
4. Analyze the generated report to gain insights into your team's productivity

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

When contributing:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/new-metric`)
3. Commit your changes (`git commit -am 'Add new metric type'`)
4. Push to the branch (`git push origin feature/new-metric`)
5. Create a new Pull Request

Please make sure to update tests as appropriate and follow the modular code structure.
