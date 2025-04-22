# Kanban Reports

A Go application for generating reports from Kanban board data exported as CSV. This tool helps track and analyze team productivity by generating reports on story points delivered, broken down by contributor, epic, product area, or team.

## Features

- Process CSV exports from kanban boards
- Generate reports by:
  - Contributor (who completed the work)
  - Epic (which larger initiatives the work belongs to)
  - Product Area (which product categories the work affects)
  - Team (which teams are delivering the most)
- Filter by date ranges or last N days
- Save reports to file or view in console

## Installation

### Prerequisites

- Go 1.21 or higher

### From Source

```zsh
# Clone the repository
git clone https://github.com/hannasdev/kanban-reports.git
cd kanban-reports

# Build the application
go build -o bin/kanban-reports ./cmd/kanban-reports
```

## Usage

### Basic Usage

```zsh
# Generate a contributor report for the last 7 days
./bin/kanban-reports --csv data/kanban-data.csv --type contributor --last 7

# Generate an epic report for a specific date range
./bin/kanban-reports --csv data/kanban-data.csv --type epic --start 2024-05-01 --end 2024-05-31

# Save the report to a file
./bin/kanban-reports --csv data/kanban-data.csv --type product-area --last 14 --output product-report.txt
```

### Command Line Options

| Flag | Description | Example |
|------|-------------|---------|
| `--csv` | Path to the kanban CSV file (required) | `--csv data/kanban-data.csv` |
| `--type` | Type of report to generate (contributor, epic, product-area, team) | `--type epic` |
| `--start` | Start date in YYYY-MM-DD format | `--start 2024-05-01` |
| `--end` | End date in YYYY-MM-DD format | `--end 2024-05-31` |
| `--last` | Generate report for the last N days | `--last 7` |
| `--output` | Path to save the report (optional) | `--output report.txt` |
| `--delimiter` | CSV delimiter: comma, tab, semicolon, or auto (default: auto) | `--delimiter comma` |

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

```zsh
id,name,type,requester,owners,description,is_completed,created_at,started_at,updated_at,moved_at,completed_at,estimate,external_ticket_count,external_tickets,is_blocked,is_a_blocker,due_date,labels,epic_labels,tasks,state,epic_id,epic,project_id,project,iteration_id,iteration,utc_offset,is_archived,team_id,team,epic_state,epic_is_archived,epic_created_at,epic_started_at,epic_due_date,milestone_id,milestone,milestone_state,milestone_created_at,milestone_started_at,milestone_due_date,milestone_categories,epic_planned_start_date,workflow,workflow_id,priority,severity,product_area,skill_set,technical_area,custom_fields
```

## Report Types

### Contributor Report

Shows story points completed by each contributor (owner). Points are divided equally if multiple owners are assigned to an item.

Example output:

```zsh
Story Points by Contributor:

john.smith@example.com          15.0 points   5 items
jane.doe@example.com            12.5 points   4 items
Unassigned                       5.0 points   2 items

Total: 32.5 points across 11 items
```

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

Total: 48.5 points across 19 item
```
