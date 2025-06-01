// internal/parser/csv_parser_test.go
package parser

import (
	"os"
	"strings"
	"testing"

	"github.com/hannasdev/kanban-reports/internal/models"
)

func TestCSVParser_Parse(t *testing.T) {
	// Create a temporary CSV file for testing
	tempFile, err := os.CreateTemp("", "test-kanban-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data to the file
	testCSV := `id,name,type,estimate,is_completed,completed_at,owners,epic,team,product_area
1,Task 1,Feature,3,TRUE,2024/05/07 10:30:00,john@example.com,Epic 1,Team A,Backend
2,Task 2,Bug,1,TRUE,2024/05/08 15:45:00,jane@example.com,Epic 1,Team A,Frontend
3,Task 3,Task,5,FALSE,,john@example.com;jane@example.com,Epic 2,Team B,Backend
`
	if _, err := tempFile.Write([]byte(testCSV)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Parse the CSV file
	parser := NewCSVParser(tempFile.Name())
	items, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Verify the results
	if len(items) != 3 {
		t.Errorf("Parse() returned %d items, want 3", len(items))
	}

	// Check the first item
	if items[0].ID != "1" || items[0].Name != "Task 1" || items[0].Type != "Feature" {
		t.Errorf("Parse() first item = %+v, want ID=1, Name=Task 1, Type=Feature", items[0])
	}

	// Check completed status
	if !items[0].IsCompleted || items[2].IsCompleted {
		t.Errorf("Parse() completed status incorrect")
	}

	// Check owners parsing
	if len(items[2].Owners) != 2 {
		t.Errorf("Parse() owners = %v, want 2 owners", items[2].Owners)
	}
}

func TestParseRow(t *testing.T) {
	// Create a test row and column indices
	row := []string{
		"TASK-123",                      // id
		"Important Task",                // name
		"Feature",                       // type
		"User",                          // requester
		"john@example.com,jane@example.com", // owners
		"Task description",              // description
		"TRUE",                          // is_completed
		"2024/05/01 10:00:00",           // created_at
		"2024/05/02 11:00:00",           // started_at
		"2024/05/07 15:00:00",           // updated_at
		"2024/05/05 12:00:00",           // moved_at
		"2024/05/07 14:00:00",           // completed_at
		"5",                             // estimate
		"2",                             // external_ticket_count
		"{\"JIRA-123\":1,\"JIRA-456\":1}", // external_tickets
		"FALSE",                         // is_blocked
		"FALSE",                         // is_a_blocker
		"2024/05/15 00:00:00",           // due_date
		"feature,priority",              // labels
		"epic-label-1,epic-label-2",     // epic_labels
		"task1,task2,task3",             // tasks
		"In Progress",                   // state
		"EPIC-456",                      // epic_id
		"Major Epic",                    // epic
		"PROJ-789",                      // project_id
		"Core Project",                  // project
		"ITER-101",                      // iteration_id
		"Sprint 42",                     // iteration
		"UTC+0",                         // utc_offset
		"FALSE",                         // is_archived
		"TEAM-202",                      // team_id
		"Engineering",                   // team
		"In Progress",                   // epic_state
		"FALSE",                         // epic_is_archived
		"2024/04/01 09:00:00",           // epic_created_at
		"2024/04/02 10:00:00",           // epic_started_at
		"2024/06/01 00:00:00",           // epic_due_date
		"MILE-303",                      // milestone_id
		"Q2 Release",                    // milestone
		"In Progress",                   // milestone_state
		"2024/04/01 00:00:00",           // milestone_created_at
		"2024/04/01 00:00:00",           // milestone_started_at
		"2024/06/30 00:00:00",           // milestone_due_date
		"release,major",                 // milestone_categories
		"2024/04/01 00:00:00",           // epic_planned_start_date
		"Kanban",                        // workflow
		"WF-404",                        // workflow_id
		"High",                          // priority
		"Medium",                        // severity
		"Backend",                       // product_area
		"Go,API",                        // skill_set
		"Databases",                     // technical_area
		"importance=high;domain=core",   // custom_fields
	}
	
	colIndices := map[string]int{
		"id": 0,
		"name": 1,
		"type": 2,
		"requester": 3,
		"owners": 4,
		"description": 5,
		"is_completed": 6,
		"created_at": 7,
		"started_at": 8,
		"updated_at": 9,
		"moved_at": 10,
		"completed_at": 11,
		"estimate": 12,
		"external_ticket_count": 13,
		"external_tickets": 14,
		"is_blocked": 15,
		"is_a_blocker": 16,
		"due_date": 17,
		"labels": 18,
		"epic_labels": 19,
		"tasks": 20,
		"state": 21,
		"epic_id": 22,
		"epic": 23,
		"project_id": 24,
		"project": 25,
		"iteration_id": 26,
		"iteration": 27,
		"utc_offset": 28,
		"is_archived": 29,
		"team_id": 30,
		"team": 31,
		"epic_state": 32,
		"epic_is_archived": 33,
		"epic_created_at": 34,
		"epic_started_at": 35,
		"epic_due_date": 36,
		"milestone_id": 37,
		"milestone": 38,
		"milestone_state": 39,
		"milestone_created_at": 40,
		"milestone_started_at": 41,
		"milestone_due_date": 42,
		"milestone_categories": 43,
		"epic_planned_start_date": 44,
		"workflow": 45,
		"workflow_id": 46,
		"priority": 47,
		"severity": 48,
		"product_area": 49,
		"skill_set": 50,
		"technical_area": 51,
		"custom_fields": 52,
	}
	
	// Create a parser
	parser := NewCSVParser("test.csv")
	
	// Parse the row
	item, err := parser.parseRow(row, colIndices)
	if err != nil {
		t.Fatalf("parseRow() error = %v", err)
	}
	
	// Verify the parsed item
	tests := []struct {
		name     string
		actual   interface{}
		expected interface{}
	}{
		{"ID", item.ID, "TASK-123"},
		{"Name", item.Name, "Important Task"},
		{"Type", item.Type, "Feature"},
		{"Requester", item.Requester, "User"},
		{"Owners count", len(item.Owners), 2},
		{"Description", item.Description, "Task description"},
		{"IsCompleted", item.IsCompleted, true},
		{"Estimate", item.Estimate, 5.0},
		{"ExternalTicketCount", item.ExternalTicketCount, 2},
		{"ExternalTickets count", len(item.ExternalTickets), 2},
		{"IsBlocked", item.IsBlocked, false},
		{"IsABlocker", item.IsABlocker, false},
		{"Labels count", len(item.Labels), 2},
		{"EpicLabels count", len(item.EpicLabels), 2},
		{"Tasks count", len(item.Tasks), 3},
		{"State", item.State, "In Progress"},
		{"EpicID", item.EpicID, "EPIC-456"},
		{"Epic", item.Epic, "Major Epic"},
		{"ProjectID", item.ProjectID, "PROJ-789"},
		{"Project", item.Project, "Core Project"},
		{"IterationID", item.IterationID, "ITER-101"},
		{"Iteration", item.Iteration, "Sprint 42"},
		{"UTCOffset", item.UTCOffset, "UTC+0"},
		{"IsArchived", item.IsArchived, false},
		{"TeamID", item.TeamID, "TEAM-202"},
		{"Team", item.Team, "Engineering"},
		{"EpicState", item.EpicState, "In Progress"},
		{"EpicIsArchived", item.EpicIsArchived, false},
		{"MilestoneID", item.MilestoneID, "MILE-303"},
		{"Milestone", item.Milestone, "Q2 Release"},
		{"MilestoneState", item.MilestoneState, "In Progress"},
		{"MilestoneCategories count", len(item.MilestoneCategories), 2},
		{"Workflow", item.Workflow, "Kanban"},
		{"WorkflowID", item.WorkflowID, "WF-404"},
		{"Priority", item.Priority, "High"},
		{"Severity", item.Severity, "Medium"},
		{"ProductArea", item.ProductArea, "Backend"},
		{"SkillSet", item.SkillSet, "Go,API"},
		{"TechnicalArea", item.TechnicalArea, "Databases"},
		{"CustomFields count", len(item.CustomFields), 2},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.actual, tt.expected)
			}
		})
	}
	
	// Verify timestamps
	createdAt, _ := models.ParseTime("2024/05/01 10:00:00")
	if !item.CreatedAt.Equal(createdAt) {
		t.Errorf("CreatedAt = %v, want %v", item.CreatedAt, createdAt)
	}
	
	completedAt, _ := models.ParseTime("2024/05/07 14:00:00")
	if !item.CompletedAt.Equal(completedAt) {
		t.Errorf("CompletedAt = %v, want %v", item.CompletedAt, completedAt)
	}
}

func TestCSVParser_FileErrors(t *testing.T) {
	tests := []struct {
		name        string
		setupFile   func() string // Returns file path
		expectError bool
		errorType   string
	}{
		{
			name: "Nonexistent file",
			setupFile: func() string {
				return "/nonexistent/path/file.csv"
			},
			expectError: true,
			errorType:   "does not exist", // Updated to match new error message
		},
		{
			name: "Empty file",
			setupFile: func() string {
				tempFile, _ := os.CreateTemp("", "empty-*.csv")
				tempFile.Close()
				return tempFile.Name()
			},
			expectError: true,
			errorType:   "error reading CSV header",
		},
		{
			name: "File with only whitespace",
			setupFile: func() string {
					tempFile, _ := os.CreateTemp("", "whitespace-*.csv")
					tempFile.WriteString("   \n  \n  ")
					tempFile.Close()
					return tempFile.Name()
			},
			expectError: true,
			errorType:   "required column", // Actual error when headers are whitespace
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setupFile()
			defer os.Remove(filePath) // Clean up

			parser := NewCSVParser(filePath)
			_, err := parser.Parse()

			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err != nil)
			}

			if tt.expectError && err != nil {
				if !strings.Contains(err.Error(), tt.errorType) {
					t.Errorf("Expected error containing '%s', got: %s", tt.errorType, err.Error())
				}
			}
		})
	}
}

func TestCSVParser_InvalidCSVStructure(t *testing.T) {
	tests := []struct {
		name        string
		csvContent  string
		expectError bool
		description string
	}{
		{
			name: "Missing required columns",
			csvContent: `name,description,owners
Task 1,Description 1,john@example.com
Task 2,Description 2,jane@example.com`,
			expectError: true,
			description: "Should fail when required columns (id, estimate, is_completed, completed_at) are missing",
		},
		{
			name: "Headers only, no data rows",
			csvContent: `id,name,estimate,is_completed,completed_at`,
			expectError: false,
			description: "Should return empty slice when no data rows present",
		},
		{
			name: "Inconsistent field counts per row",
			csvContent: `id,name,estimate,is_completed,completed_at
1,Task 1,3,TRUE,2024/05/01 10:00:00
2,Task 2,2,FALSE
3,Task 3,5,TRUE,2024/05/03 10:00:00,extra,field`,
			expectError: false,
			description: "Should handle inconsistent field counts gracefully (CSV parser allows this)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file with test content
			tempFile, err := os.CreateTemp("", "csv-test-*.csv")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tempFile.Name())

			if _, err := tempFile.WriteString(tt.csvContent); err != nil {
				t.Fatalf("Failed to write test content: %v", err)
			}
			tempFile.Close()

			// Test parsing
			parser := NewCSVParser(tempFile.Name())
			items, err := parser.Parse()

			if (err != nil) != tt.expectError {
				t.Errorf("%s: Expected error: %v, got: %v", tt.description, tt.expectError, err != nil)
			}

			if !tt.expectError && len(items) == 0 && tt.name != "Headers only, no data rows" {
				t.Errorf("%s: Expected items but got empty slice", tt.description)
			}
		})
	}
}

func TestCSVParser_DataTypeHandling(t *testing.T) {
	tests := []struct {
		name       string
		csvContent string
		itemIndex  int
		validate   func(models.KanbanItem) bool
	}{
		{
			name: "Invalid date formats",
			csvContent: `id,name,estimate,is_completed,completed_at,created_at,started_at
1,Task 1,3,TRUE,invalid-date,2024/05/01 10:00:00,2024/05/02 10:00:00
2,Task 2,2,FALSE,,2024/05/01 10:00:00,`,
			itemIndex: 0,
			validate: func(item models.KanbanItem) bool {
				// Invalid date should result in zero time
				return item.CompletedAt.IsZero() && !item.CreatedAt.IsZero()
			},
		},
		{
			name: "Invalid boolean values",
			csvContent: `id,name,estimate,is_completed,completed_at,is_blocked,is_a_blocker
1,Task 1,3,maybe,2024/05/01 10:00:00,yes,no
2,Task 2,2,1,2024/05/02 10:00:00,0,1`,
			itemIndex: 0,
			validate: func(item models.KanbanItem) bool {
				// Invalid boolean should default to false
				return !item.IsCompleted && !item.IsBlocked && !item.IsABlocker
			},
		},
		{
			name: "Invalid numeric values",
			csvContent: `id,name,estimate,is_completed,completed_at,external_ticket_count
1,Task 1,not-a-number,TRUE,2024/05/01 10:00:00,also-not-a-number
2,Task 2,,FALSE,2024/05/02 10:00:00,`,
			itemIndex: 0,
			validate: func(item models.KanbanItem) bool {
				// Invalid numbers should default to 0
				return item.Estimate == 0 && item.ExternalTicketCount == 0
			},
		},
		{
			name: "Complex external tickets JSON",
			csvContent: `id,name,estimate,is_completed,completed_at,external_tickets
	1,Task 1,3,TRUE,2024/05/01 10:00:00,"#{""JIRA-123"":1,""GITHUB-456"":1}"
	2,Task 2,2,FALSE,2024/05/02 10:00:00,invalid-json`,
			itemIndex: 0,
			validate: func(item models.KanbanItem) bool {
					// Should parse valid JSON tickets
					return len(item.ExternalTickets) == 2
			},
		},
		{
			name: "Various owner formats",
			csvContent: `id,name,estimate,is_completed,completed_at,owners
1,Task 1,3,TRUE,2024/05/01 10:00:00,"john@example.com, jane@example.com"
2,Task 2,2,FALSE,2024/05/02 10:00:00,bob@example.com;alice@example.com
3,Task 3,1,TRUE,2024/05/03 10:00:00,   spaced@example.com   `,
			itemIndex: 0,
			validate: func(item models.KanbanItem) bool {
				// Should handle comma-separated owners correctly
				return len(item.Owners) >= 2
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file with test content
			tempFile, err := os.CreateTemp("", "csv-datatype-*.csv")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tempFile.Name())

			if _, err := tempFile.WriteString(tt.csvContent); err != nil {
				t.Fatalf("Failed to write test content: %v", err)
			}
			tempFile.Close()

			// Test parsing
			parser := NewCSVParser(tempFile.Name())
			items, err := parser.Parse()

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(items) <= tt.itemIndex {
				t.Fatalf("Expected at least %d items, got %d", tt.itemIndex+1, len(items))
			}

			if !tt.validate(items[tt.itemIndex]) {
				t.Errorf("Validation failed for item at index %d", tt.itemIndex)
			}
		})
	}
}

func TestCSVParser_DelimiterHandling(t *testing.T) {
	tests := []struct {
		name       string
		csvContent string
		delimiter  models.DelimiterType
		expectRows int
	}{
		{
			name:       "Tab delimited with auto detection",
			csvContent: "id\tname\testimate\tis_completed\tcompleted_at\n1\tTask 1\t3\tTRUE\t2024/05/01 10:00:00",
			delimiter:  models.DelimiterAuto,
			expectRows: 1,
		},
		{
			name:       "Semicolon delimited with explicit setting",
			csvContent: "id;name;estimate;is_completed;completed_at\n1;Task 1;3;TRUE;2024/05/01 10:00:00",
			delimiter:  models.DelimiterSemicolon,
			expectRows: 1,
		},
		{
			name:       "Comma delimited content",
			csvContent: "id,name,estimate,is_completed,completed_at\n1,Task 1,3,TRUE,2024/05/01 10:00:00",
			delimiter:  models.DelimiterAuto,
			expectRows: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file with test content
			tempFile, err := os.CreateTemp("", "csv-delimiter-*.csv")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tempFile.Name())

			if _, err := tempFile.WriteString(tt.csvContent); err != nil {
				t.Fatalf("Failed to write test content: %v", err)
			}
			tempFile.Close()

			// Test parsing with specified delimiter
			parser := NewCSVParser(tempFile.Name()).WithDelimiter(tt.delimiter)
			items, err := parser.Parse()

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(items) != tt.expectRows {
				t.Errorf("Expected %d rows, got %d", tt.expectRows, len(items))
			}
		})
	}
}

func TestCSVParser_ErrorRecovery(t *testing.T) {
	// Test that parser continues processing even when some rows have errors
	csvContent := `id,name,estimate,is_completed,completed_at
1,Task 1,3,TRUE,2024/05/01 10:00:00
,Task 2,2,FALSE,2024/05/02 10:00:00
3,Task 3,invalid-estimate,TRUE,invalid-date
4,Task 4,1,TRUE,2024/05/04 10:00:00`

	tempFile, err := os.CreateTemp("", "csv-recovery-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write test content: %v", err)
	}
	tempFile.Close()

	parser := NewCSVParser(tempFile.Name())
	items, err := parser.Parse()

	if err != nil {
		t.Fatalf("Expected parser to recover from row errors, got: %v", err)
	}

	// Should have successfully parsed some rows despite errors in others
	if len(items) == 0 {
		t.Errorf("Expected some successfully parsed items despite row errors")
	}

	// Verify that valid rows were parsed correctly
	validItems := 0
	for _, item := range items {
		if item.ID != "" && item.Name != "" {
			validItems++
		}
	}

	if validItems == 0 {
		t.Errorf("Expected at least some valid items to be parsed")
	}
}