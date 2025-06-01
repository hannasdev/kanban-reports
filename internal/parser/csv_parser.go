package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
)

const (
	// DelimiterDetectionBufferSize is the buffer size for automatic delimiter detection
	DelimiterDetectionBufferSize = 4 * 1024 // 4KB
)

var (
	// RequiredColumns are the minimum columns needed for parsing
	RequiredColumns = []string{"id", "name", "estimate", "is_completed", "completed_at"}
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
	file, err := p.openAndPrepareFile()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := p.createCSVReader(file)
	
	_, colIndices, err := p.parseHeaders(reader)
	if err != nil {
		return nil, err
	}

	if err := p.validateRequiredColumns(colIndices); err != nil {
		return nil, err
	}

	items, err := p.parseDataRows(reader, colIndices)
	if err != nil {
		return nil, err
	}

	fmt.Printf("✅ Loaded %d kanban items\n", len(items))
	return items, nil
}

// openAndPrepareFile opens the CSV file and handles delimiter detection
func (p *CSVParser) openAndPrepareFile() (*os.File, error) {
	file, err := os.Open(p.filepath)
	if err != nil {
		return nil, p.formatFileError(err)
	}

	// Handle automatic delimiter detection
	if p.delimiter.AutoDetect {
		if err := p.detectDelimiter(file); err != nil {
			file.Close()
			return nil, err
		}
		// Reset file pointer after detection
		if _, err := file.Seek(0, 0); err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to reset file pointer: %w", err)
		}
	}

	return file, nil
}

// formatFileError provides user-friendly file access error messages
func (p *CSVParser) formatFileError(err error) error {
	if os.IsNotExist(err) {
		return fmt.Errorf("CSV file '%s' does not exist", p.filepath)
	}
	if os.IsPermission(err) {
		return fmt.Errorf("permission denied accessing CSV file '%s'", p.filepath)
	}
	
	// Check if it's a directory
	if info, statErr := os.Stat(p.filepath); statErr == nil && info.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a file. Please specify a CSV file path", p.filepath)
	}
	
	return fmt.Errorf("error opening CSV file '%s': %w", p.filepath, err)
}

// detectDelimiter reads a sample of the file to detect the CSV delimiter
func (p *CSVParser) detectDelimiter(file *os.File) error {
	buffer := make([]byte, DelimiterDetectionBufferSize)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file for delimiter detection: %w", err)
	}
	
	sampleContent := string(buffer[:n])
	p.delimiter = models.DetectDelimiterType(sampleContent)
	fmt.Printf("Detected %s-delimited CSV\n", p.delimiter.Name)
	
	return nil
}

// createCSVReader creates and configures a CSV reader
func (p *CSVParser) createCSVReader(file *os.File) *csv.Reader {
	reader := csv.NewReader(file)
	reader.Comma = p.delimiter.Value
	reader.FieldsPerRecord = -1 // Disable field count checking for flexibility
	return reader
}

// parseHeaders reads and processes the CSV header row
func (p *CSVParser) parseHeaders(reader *csv.Reader) ([]string, map[string]int, error) {
	headers, err := reader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("error reading CSV header: %w", err)
	}

	// Create column index map for fast lookup
	colIndices := make(map[string]int)
	for i, header := range headers {
		colIndices[strings.TrimSpace(header)] = i
	}

	fmt.Println("Found columns:", strings.Join(headers, ", "))
	return headers, colIndices, nil
}

// validateRequiredColumns ensures all required columns are present
func (p *CSVParser) validateRequiredColumns(colIndices map[string]int) error {
	for _, col := range RequiredColumns {
		if _, exists := colIndices[col]; !exists {
			return fmt.Errorf("required column '%s' not found in CSV headers", col)
		}
	}
	return nil
}

// parseDataRows reads and parses all data rows from the CSV
func (p *CSVParser) parseDataRows(reader *csv.Reader, colIndices map[string]int) ([]models.KanbanItem, error) {
	var items []models.KanbanItem
	rowNumber := 1 // Start at 1 since we already read the header
	
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV row %d: %w", rowNumber, err)
		}

		item, err := p.parseRow(row, colIndices)
		if err != nil {
			// Log warning but continue processing
			fmt.Printf("Warning: error parsing row %d: %v\n", rowNumber, err)
			rowNumber++
			continue
		}

		items = append(items, item)
		rowNumber++
	}

	return items, nil
}

// parseRow converts a CSV row into a KanbanItem
func (p *CSVParser) parseRow(row []string, colIndices map[string]int) (models.KanbanItem, error) {
	item := models.KanbanItem{}
	
	// Helper function to safely get column values
	getCol := func(name string) string {
		if idx, exists := colIndices[name]; exists && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	// Parse basic fields
	if err := p.parseBasicFields(&item, getCol); err != nil {
		return item, fmt.Errorf("failed to parse basic fields: %w", err)
	}

	// Parse timestamps
	if err := p.parseTimestamps(&item, getCol); err != nil {
		return item, fmt.Errorf("failed to parse timestamps: %w", err)
	}

	// Parse numeric fields
	if err := p.parseNumericFields(&item, getCol); err != nil {
		return item, fmt.Errorf("failed to parse numeric fields: %w", err)
	}

	// Parse collection fields (arrays, maps)
	p.parseCollectionFields(&item, getCol)

	// Parse organizational fields
	p.parseOrganizationalFields(&item, getCol)

	return item, nil
}

// parseBasicFields sets basic string fields on the KanbanItem
func (p *CSVParser) parseBasicFields(item *models.KanbanItem, getCol func(string) string) error {
	item.ID = getCol("id")
	item.Name = getCol("name")
	item.Type = getCol("type")
	item.Requester = getCol("requester")
	item.Description = getCol("description")
	item.State = getCol("state")
	item.UTCOffset = getCol("utc_offset")
	item.Workflow = getCol("workflow")
	item.WorkflowID = getCol("workflow_id")
	item.Priority = getCol("priority")
	item.Severity = getCol("severity")
	item.ProductArea = getCol("product_area")
	item.SkillSet = getCol("skill_set")
	item.TechnicalArea = getCol("technical_area")

	// Validate required fields
	if item.ID == "" {
		return fmt.Errorf("missing required field: id")
	}
	if item.Name == "" {
		return fmt.Errorf("missing required field: name")
	}

	return nil
}

// parseTimestamps parses all timestamp fields
func (p *CSVParser) parseTimestamps(item *models.KanbanItem, getCol func(string) string) error {
	timestampFields := map[string]*time.Time{
		"created_at":                &item.CreatedAt,
		"started_at":                &item.StartedAt,
		"updated_at":                &item.UpdatedAt,
		"moved_at":                  &item.MovedAt,
		"completed_at":              &item.CompletedAt,
		"due_date":                  &item.DueDate,
		"epic_created_at":           &item.EpicCreatedAt,
		"epic_started_at":           &item.EpicStartedAt,
		"epic_due_date":             &item.EpicDueDate,
		"milestone_created_at":      &item.MilestoneCreatedAt,
		"milestone_started_at":      &item.MilestoneStartedAt,
		"milestone_due_date":        &item.MilestoneDueDate,
		"epic_planned_start_date":   &item.EpicPlannedStartDate,
	}

	for fieldName, timePtr := range timestampFields {
		if timeStr := getCol(fieldName); timeStr != "" {
			if parsedTime, err := models.ParseTime(timeStr); err == nil {
				*timePtr = parsedTime
			}
			// Ignore parse errors for optional timestamp fields
		}
	}

	return nil
}

// parseNumericFields parses numeric fields with validation
func (p *CSVParser) parseNumericFields(item *models.KanbanItem, getCol func(string) string) error {
	// Parse boolean fields
	item.IsCompleted = models.ParseBool(getCol("is_completed"))
	item.IsBlocked = models.ParseBool(getCol("is_blocked"))
	item.IsABlocker = models.ParseBool(getCol("is_a_blocker"))
	item.IsArchived = models.ParseBool(getCol("is_archived"))
	item.EpicIsArchived = models.ParseBool(getCol("epic_is_archived"))

	// Parse numeric fields
	item.Estimate = models.ParseFloat(getCol("estimate"))
	item.ExternalTicketCount = models.ParseInt(getCol("external_ticket_count"))

	return nil
}

// parseCollectionFields parses array and map fields
func (p *CSVParser) parseCollectionFields(item *models.KanbanItem, getCol func(string) string) {
	item.Owners = models.ParseOwners(getCol("owners"))
	item.Labels = models.ParseStringList(getCol("labels"))
	item.EpicLabels = models.ParseStringList(getCol("epic_labels"))
	item.Tasks = models.ParseStringList(getCol("tasks"))
	item.ExternalTickets = models.ParseExternalTickets(getCol("external_tickets"))
	item.MilestoneCategories = models.ParseStringList(getCol("milestone_categories"))
	item.CustomFields = models.ParseCustomFields(getCol("custom_fields"))
}

// parseOrganizationalFields parses project, epic, team, and milestone fields
func (p *CSVParser) parseOrganizationalFields(item *models.KanbanItem, getCol func(string) string) {
	// Epic fields
	item.EpicID = getCol("epic_id")
	item.Epic = getCol("epic")
	item.EpicState = getCol("epic_state")

	// Project fields
	item.ProjectID = getCol("project_id")
	item.Project = getCol("project")

	// Iteration fields
	item.IterationID = getCol("iteration_id")
	item.Iteration = getCol("iteration")

	// Team fields
	item.TeamID = getCol("team_id")
	item.Team = getCol("team")

	// Milestone fields
	item.MilestoneID = getCol("milestone_id")
	item.Milestone = getCol("milestone")
	item.MilestoneState = getCol("milestone_state")
}