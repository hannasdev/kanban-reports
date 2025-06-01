# Kanban Reports

A Go application for generating productivity reports and metrics from Kanban board data exported as CSV. Track team performance, analyze flow efficiency, and identify improvement opportunities with comprehensive reporting and metrics analysis.

## 🚀 Quick Start

### Option 1: Interactive Mode (Recommended for beginners)

```bash
# Build and run setup
./scripts/setup.sh

# Start interactive mode
./bin/kanban-reports --interactive
```

### Option 2: Command Line Mode

```bash
# Show all available options
./bin/kanban-reports --help

# See practical examples
./bin/kanban-reports --examples

# Quick example with sample data
./bin/kanban-reports --csv data/sample.csv --type contributor --last 30
```

## 📊 Features

### Reports

- **Contributor Reports**: Story points by person who completed work
- **Epic Reports**: Story points by epic/initiative
- **Product Area Reports**: Story points by product category
- **Team Reports**: Story points by team

### Advanced Metrics

- **Lead Time Analysis**: How long items take from creation to completion
- **Throughput Analysis**: Completion rates over time (items & points)
- **Flow Efficiency**: Active vs waiting time analysis
- **Estimation Accuracy**: Correlation between estimates and actual time
- **Work Item Age**: Age analysis of current incomplete work
- **Team Improvement**: Month-over-month improvement trends

### Filtering & Output

- Filter by date ranges or last N days
- Filter ad-hoc requests (include, exclude, or focus only on them)
- Save reports to file or view in console
- Automatic CSV delimiter detection (comma, tab, semicolon)
- Multiple date field filtering options

## 🔗 Shortcut.com Integration

This application is specifically designed to work with CSV exports from [Shortcut.com](https://shortcut.com) (formerly Clubhouse).

### Exporting Data from Shortcut.com

1. **Navigate to Stories**: Go to your Shortcut workspace and click on "Stories"
2. **Apply Filters**: Use Shortcut's search and filter options to select the stories you want to analyze
3. **Export to CSV**:
   - Click the "Export" button in the top right
   - Select "CSV" format
   - Choose "All fields" for complete data export
4. **Save the File**: Download the CSV file to your local machine

### Supported Shortcut Fields

The application automatically maps Shortcut.com CSV columns including:

| Shortcut Field | Used For | Report Types |
|----------------|----------|--------------|
| `owners` | Contributor attribution | Contributor reports |
| `epic` | Epic grouping | Epic reports |
| `team` | Team attribution | Team reports |
| `estimate` | Story points | All reports & metrics |
| `completed_at` | Completion tracking | Date filtering, metrics |
| `created_at` | Lead time calculation | Lead time metrics |
| `started_at` | Cycle time calculation | Flow & cycle time metrics |
| `labels` | Ad-hoc filtering | Filtering (looks for "ad-hoc-request" label) |
| `product_area` | Product categorization | Product area reports |

### Tips for Shortcut Users

- **Use Labels**: Tag ad-hoc requests with "ad-hoc-request" label for filtering
- **Set Story Points**: Ensure stories have estimates for meaningful reports
- **Track State Changes**: Use Shortcut's workflow states for accurate timing data
- **Regular Exports**: Export data weekly/monthly for trend analysis

## 📦 Installation

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

### Using Make (Recommended)

```bash
# Build, test, and set up everything
make all

# Or just build
make build

# Run tests
make test

# Install to GOPATH/bin
make install
```

## 📖 Usage Examples

### Basic Reports

```bash
# Team productivity this month
./bin/kanban-reports --csv kanban-data.csv --type contributor --last 30

# Epic progress over the quarter
./bin/kanban-reports --csv kanban-data.csv --type epic --last 90 --output epic-report.txt

# Product area breakdown for specific period
./bin/kanban-reports --csv kanban-data.csv --type product-area --start 2024-01-01 --end 2024-03-31
```

### Metrics Analysis

```bash
# Analyze lead times by story point size
./bin/kanban-reports --csv kanban-data.csv --metrics lead-time --last 90

# Weekly throughput trends
./bin/kanban-reports --csv kanban-data.csv --metrics throughput --period week --last 180

# Complete metrics analysis
./bin/kanban-reports --csv kanban-data.csv --metrics all --last 90 --output full-analysis.txt
```

### Advanced Filtering

```bash
# Exclude ad-hoc work to see planned work only
./bin/kanban-reports --csv kanban-data.csv --type team --last 30 --ad-hoc exclude

# Analyze only urgent/ad-hoc requests
./bin/kanban-reports --csv kanban-data.csv --metrics throughput --ad-hoc only --last 60

# Filter by creation date instead of completion date
./bin/kanban-reports --csv kanban-data.csv --type contributor --last 30 --filter-field created_at
```

## ⚙️ Command Line Options

| Flag | Description | Example |
|------|-------------|---------|
| `--help, -h` | Complete help with all options | `./bin/kanban-reports --help` |
| `--examples` | Practical usage examples | `./bin/kanban-reports --examples` |
| `--version` | Version information | `./bin/kanban-reports --version` |
| `--interactive, -i` | Interactive menu mode | `./bin/kanban-reports -i` |
| `--csv` | Path to the kanban CSV file (required) | `--csv data/kanban-data.csv` |
| `--type` | Report type (contributor, epic, product-area, team) | `--type epic` |
| `--metrics` | Metrics type (lead-time, throughput, flow, estimation, age, improvement, all) | `--metrics lead-time` |
| `--period` | Time period for metrics (week, month) | `--period week` |
| `--start` | Start date (YYYY-MM-DD) | `--start 2024-05-01` |
| `--end` | End date (YYYY-MM-DD) | `--end 2024-05-31` |
| `--last` | Last N days | `--last 7` |
| `--output` | Save to file | `--output report.txt` |
| `--delimiter` | CSV delimiter (comma, tab, semicolon, auto) | `--delimiter comma` |
| `--ad-hoc` | Ad-hoc filter (include, exclude, only) | `--ad-hoc exclude` |

## 📋 CSV Data Format

### Required Columns

- `id`: Unique identifier for the item
- `name`: Name or title of the item
- `estimate`: Story points or estimate value
- `is_completed`: Boolean indicating if the item is completed
- `completed_at`: Date when the item was completed

### Optional but Recommended Columns

- `owners`: Person(s) assigned to the item (semicolon-separated)
- `epic`: Epic name
- `team`: Team name
- `product_area`: Product area or category
- `type`: Type of work (feature, bug, task, etc.)
- `created_at`: When the item was created
- `started_at`: When work began on the item
- `labels`: Labels (use "ad-hoc-request" for filtering)

### Example CSV Header

```csv
id,name,type,requester,owners,description,is_completed,created_at,started_at,updated_at,moved_at,completed_at,estimate,external_ticket_count,external_tickets,is_blocked,is_a_blocker,due_date,labels,epic_labels,tasks,state,epic_id,epic,project_id,project,iteration_id,iteration,utc_offset,is_archived,team_id,team,epic_state,epic_is_archived,epic_created_at,epic_started_at,epic_due_date,milestone_id,milestone,milestone_state,milestone_created_at,milestone_started_at,milestone_due_date,milestone_categories,epic_planned_start_date,workflow,workflow_id,priority,severity,product_area,skill_set,technical_area,custom_fields
```

## 📊 Sample Output

### Contributor Report

```txt
Story Points by Contributor:

john@example.com                 25.0 points   8 items
jane@example.com                 18.5 points   6 items
bob@example.com                  12.0 points   4 items
Unassigned                        3.0 points   2 items

Total: 58.5 points across 20 items
```

### Lead Time Metrics

```txt
# Lead Time Analysis by Story Point Size (in days)

Story points | Count | Min | Max | Avg | Median
-------------|-------|-----|-----|-----|-------
           1 |    15 | 2.5 | 8.3 | 4.2 |    3.8
           3 |    12 | 3.1 | 9.7 | 5.8 |    5.2
           5 |     8 | 5.7 | 15.2 | 10.1 |    9.5
           8 |     5 | 7.2 | 21.5 | 14.3 |   12.8
```

### Flow Efficiency Metrics

```txt
# Flow Efficiency Analysis

State   | Avg Time (days) | % of Total Time
--------|-----------------|---------------
Waiting |            12.5 |         71.4%
Active  |             5.0 |         28.6%

Flow Efficiency: 28.6%
```

## 🏗️ Project Structure

```txt
kanban-reports/
├── cmd/
│   └── kanban-reports/         # Main application entry point
├── internal/
│   ├── config/                 # Application configuration & CLI parsing
│   ├── menu/                   # Interactive menu system
│   ├── models/                 # Data models and types
│   ├── parser/                 # CSV parsing logic
│   ├── reports/                # Report generation
│   ├── metrics/                # Advanced metrics generation
│   └── validation/             # Input validation utilities
├── pkg/
│   ├── dateutil/               # Date handling utilities
│   ├── filtering/              # Data filtering utilities
│   └── types/                  # Shared type definitions
├── scripts/                    # Build and setup scripts
├── data/                       # Place your CSV files here
├── Makefile                    # Build automation
└── README.md
```

## 🧪 Development

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run tests for specific package
go test ./internal/reports/...

# Run with verbose output
go test -v ./internal/...
```

### Adding New Report Types

1. Add new constant in `internal/reports/types.go`
2. Create implementation file in `internal/reports/`
3. Update the switch statement in `GenerateReport`
4. Add CLI option support in `cmd/kanban-reports/main.go`

### Adding New Metrics

1. Add new constant in `internal/metrics/types.go`
2. Create implementation file in `internal/metrics/`
3. Update the switch statement in `Generator.Generate()`
4. Add CLI option support

## 🔧 Troubleshooting

### Common Issues

**"required column 'X' not found"**
- Ensure your CSV has the required columns: id, name, estimate, is_completed, completed_at
- Check that column names match exactly (case-sensitive)

**"CSV file validation failed"**
- Verify the file path is correct
- Ensure the file is readable and not corrupted
- Try using `--delimiter auto` for automatic detection

**"No items completed in the specified date range"**
- Check your date format (should be YYYY/MM/DD HH:MM:SS in CSV)
- Verify date range includes completed items
- Try using `--filter-field created_at` for broader results

### Getting Help

```bash
# Show detailed help
./bin/kanban-reports --help

# See practical examples
./bin/kanban-reports --examples

# Use interactive mode for guided setup
./bin/kanban-reports --interactive
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/new-metric`)
3. Write tests for your changes
4. Ensure tests pass (`make test`)
5. Commit your changes (`git commit -am 'Add new metric type'`)
6. Push to the branch (`git push origin feature/new-metric`)
7. Create a Pull Request

### Code Standards

- Follow Go conventions and idioms
- Add tests for new functionality
- Update documentation as needed
- Use the existing code structure and patterns

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

## 🎯 Next Steps

1. **First Time**: Run `./scripts/setup.sh` to get started quickly
2. **Export Data**: Get your CSV from Shortcut.com (or other kanban tool)
3. **Try Interactive Mode**: `./bin/kanban-reports --interactive`
4. **Explore Metrics**: Start with `--metrics all` for comprehensive analysis
5. **Automate Reports**: Set up regular exports and automated report generation

For questions, issues, or feature requests, please open an issue on GitHub.
