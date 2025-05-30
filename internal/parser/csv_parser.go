package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hannasdev/kanban-reports/internal/models"
)

// CSVParser handles parsing of kanban CSV data
type CSVParser struct {
	filepath  string
	delimiter models.DelimiterType
}

// NewCSVParser creates a new CSV parser for the specified file
func NewCSVParser(filepath string) *CSVParser {
	return &CSVParser{
		filepath:  filepath,
		delimiter: models.DelimiterComma, // Default to comma delimiter
	}
}

// WithDelimiter sets a custom delimiter for the CSV parser
func (p *CSVParser) WithDelimiter(delimiter models.DelimiterType) *CSVParser {
	p.delimiter = delimiter
	return p
}

// Parse reads the CSV file and returns a slice of KanbanItem
func (p *CSVParser) Parse() ([]models.KanbanItem, error) {
	file, err := os.Open(p.filepath)
	if err != nil {
			return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()
	
	// If auto-detection is enabled, read sample content and detect delimiter
	if p.delimiter.AutoDetect {
			// Read a sample of the file for delimiter detection
			buffer := make([]byte, 4096) // Read up to 4KB for delimiter detection
			n, _ := file.Read(buffer)
			sampleContent := string(buffer[:n])
			
			p.delimiter = models.DetectDelimiterType(sampleContent)
			fmt.Printf("Detected %s-delimited CSV\n", p.delimiter.Name)
			
			// Reset file pointer to beginning
			file.Seek(0, 0)
	}
	
	reader := csv.NewReader(file)
	
	// Set delimiter based on detection or user configuration
	reader.Comma = p.delimiter.Value
	// Disable field count checking as CSV might have inconsistent fields
	reader.FieldsPerRecord = -1
	
	// Read header row to get column indices
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV header: %w", err)
	}

	// Create a map of column name to index for easy lookup
	colIndices := make(map[string]int)
	for i, header := range headers {
		colIndices[strings.TrimSpace(header)] = i
	}

	// Debug information about found columns
	fmt.Println("Found columns:", strings.Join(headers, ", "))
	
	// Check that required columns exist
	requiredCols := []string{"id", "name", "estimate", "is_completed", "completed_at"}
	for _, col := range requiredCols {
		if _, exists := colIndices[col]; !exists {
			return nil, fmt.Errorf("required column '%s' not found in CSV headers", col)
		}
	}

	// Parse rows into KanbanItems
	var items []models.KanbanItem
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV row: %w", err)
		}

		// Create a new item and populate it
		item, err := p.parseRow(row, colIndices)
		if err != nil {
			// Log error but continue with next row
			fmt.Printf("Warning: error parsing row: %v\n", err)
			continue
		}

		items = append(items, item)
	}

	return items, nil
}

// parseRow converts a CSV row into a KanbanItem
func (p *CSVParser) parseRow(row []string, colIndices map[string]int) (models.KanbanItem, error) {
	getCol := func(name string) string {
		if idx, exists := colIndices[name]; exists && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	// Parse timestamps
	createdAt, _ := models.ParseTime(getCol("created_at"))
	startedAt, _ := models.ParseTime(getCol("started_at"))
	updatedAt, _ := models.ParseTime(getCol("updated_at"))
	movedAt, _ := models.ParseTime(getCol("moved_at"))
	completedAt, _ := models.ParseTime(getCol("completed_at"))
	dueDate, _ := models.ParseTime(getCol("due_date"))
	epicCreatedAt, _ := models.ParseTime(getCol("epic_created_at"))
	epicStartedAt, _ := models.ParseTime(getCol("epic_started_at"))
	epicDueDate, _ := models.ParseTime(getCol("epic_due_date"))
	milestoneCreatedAt, _ := models.ParseTime(getCol("milestone_created_at"))
	milestoneStartedAt, _ := models.ParseTime(getCol("milestone_started_at"))
	milestoneDueDate, _ := models.ParseTime(getCol("milestone_due_date"))
	epicPlannedStartDate, _ := models.ParseTime(getCol("epic_planned_start_date"))

	// Create the KanbanItem
	item := models.KanbanItem{
		ID:                   getCol("id"),
		Name:                 getCol("name"),
		Type:                 getCol("type"),
		Requester:            getCol("requester"),
		Owners:               models.ParseOwners(getCol("owners")),
		Description:          getCol("description"),
		IsCompleted:          models.ParseBool(getCol("is_completed")),
		CreatedAt:            createdAt,
		StartedAt:            startedAt,
		UpdatedAt:            updatedAt,
		MovedAt:              movedAt,
		CompletedAt:          completedAt,
		Estimate:             models.ParseFloat(getCol("estimate")),
		ExternalTicketCount:  models.ParseInt(getCol("external_ticket_count")),
		ExternalTickets:      models.ParseExternalTickets(getCol("external_tickets")),
		IsBlocked:            models.ParseBool(getCol("is_blocked")),
		IsABlocker:           models.ParseBool(getCol("is_a_blocker")),
		DueDate:              dueDate,
		Labels:               models.ParseStringList(getCol("labels")),
		EpicLabels:           models.ParseStringList(getCol("epic_labels")),
		Tasks:                models.ParseStringList(getCol("tasks")),
		State:                getCol("state"),
		EpicID:               getCol("epic_id"),
		Epic:                 getCol("epic"),
		ProjectID:            getCol("project_id"),
		Project:              getCol("project"),
		IterationID:          getCol("iteration_id"),
		Iteration:            getCol("iteration"),
		UTCOffset:            getCol("utc_offset"),
		IsArchived:           models.ParseBool(getCol("is_archived")),
		TeamID:               getCol("team_id"),
		Team:                 getCol("team"),
		EpicState:            getCol("epic_state"),
		EpicIsArchived:       models.ParseBool(getCol("epic_is_archived")),
		EpicCreatedAt:        epicCreatedAt,
		EpicStartedAt:        epicStartedAt,
		EpicDueDate:          epicDueDate,
		MilestoneID:          getCol("milestone_id"),
		Milestone:            getCol("milestone"),
		MilestoneState:       getCol("milestone_state"),
		MilestoneCreatedAt:   milestoneCreatedAt,
		MilestoneStartedAt:   milestoneStartedAt,
		MilestoneDueDate:     milestoneDueDate,
		MilestoneCategories:  models.ParseStringList(getCol("milestone_categories")),
		EpicPlannedStartDate: epicPlannedStartDate,
		Workflow:             getCol("workflow"),
		WorkflowID:           getCol("workflow_id"),
		Priority:             getCol("priority"),
		Severity:             getCol("severity"),
		ProductArea:          getCol("product_area"),
		SkillSet:             getCol("skill_set"),
		TechnicalArea:        getCol("technical_area"),
		CustomFields:         models.ParseCustomFields(getCol("custom_fields")),
	}

	return item, nil
}