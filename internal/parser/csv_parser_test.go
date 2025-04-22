// internal/parser/csv_parser_test.go
package parser

import (
	"os"
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