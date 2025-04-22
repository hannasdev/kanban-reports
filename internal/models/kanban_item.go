package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// KanbanItem represents a single item from the kanban board
type KanbanItem struct {
	ID                   string
	Name                 string
	Type                 string
	Requester            string
	Owners               []string
	Description          string
	IsCompleted          bool
	CreatedAt            time.Time
	StartedAt            time.Time
	UpdatedAt            time.Time
	MovedAt              time.Time
	CompletedAt          time.Time
	Estimate             float64 // Story points
	ExternalTicketCount  int
	ExternalTickets      []string
	IsBlocked            bool
	IsABlocker           bool
	DueDate              time.Time
	Labels               []string
	EpicLabels           []string
	Tasks                []string
	State                string
	EpicID               string
	Epic                 string
	ProjectID            string
	Project              string
	IterationID          string
	Iteration            string
	UTCOffset            string
	IsArchived           bool
	TeamID               string
	Team                 string
	EpicState            string
	EpicIsArchived       bool
	EpicCreatedAt        time.Time
	EpicStartedAt        time.Time
	EpicDueDate          time.Time
	MilestoneID          string
	Milestone            string
	MilestoneState       string
	MilestoneCreatedAt   time.Time
	MilestoneStartedAt   time.Time
	MilestoneDueDate     time.Time
	MilestoneCategories  []string
	EpicPlannedStartDate time.Time
	Workflow             string
	WorkflowID           string
	Priority             string
	Severity             string
	ProductArea          string
	SkillSet             string
	TechnicalArea        string
	CustomFields         map[string]string
}

// ParseTime attempts to parse time in the format provided by the CSV
func ParseTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, nil
	}
	
	// Format used in the CSV: "2024/05/07 03:49:34"
	return time.Parse("2006/01/02 15:04:05", timeStr)
}

// ParseBool converts string to bool, handling empty strings
func ParseBool(boolStr string) bool {
	if strings.ToUpper(boolStr) == "TRUE" {
		return true
	}
	return false
}

// ParseFloat converts string to float64, handling empty strings
func ParseFloat(floatStr string) float64 {
	if floatStr == "" {
		return 0
	}
	val, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		return 0
	}
	return val
}

// ParseInt converts string to int, handling empty strings
func ParseInt(intStr string) int {
	if intStr == "" {
		return 0
	}
	val, err := strconv.Atoi(intStr)
	if err != nil {
		return 0
	}
	return val
}

// ParseStringList takes a comma-separated string and returns a slice of strings
func ParseStringList(listStr string) []string {
	if listStr == "" {
		return []string{}
	}
	return strings.Split(listStr, ",")
}

// ParseExternalTickets processes the JSON-like string of external tickets
func ParseExternalTickets(ticketsStr string) []string {
	// Remove the "#" prefix if present
	if strings.HasPrefix(ticketsStr, "#") {
		ticketsStr = ticketsStr[1:]
	}
	
	// If empty or not JSON format, return empty slice
	if ticketsStr == "" || !strings.HasPrefix(ticketsStr, "{") {
		return []string{}
	}
	
	var tickets map[string]interface{}
	err := json.Unmarshal([]byte(ticketsStr), &tickets)
	if err != nil {
		return []string{}
	}
	
	// Extract keys as ticket IDs/URLs
	result := make([]string, 0, len(tickets))
	for key := range tickets {
		result = append(result, key)
	}
	
	return result
}

// ParseOwners splits owner string into individual owners
func ParseOwners(ownersStr string) []string {
	if ownersStr == "" {
		return []string{}
	}
	
	// Split by comma, semicolon, or space as potential separators
	separators := []string{",", ";", " "}
	var owners []string
	
	for _, sep := range separators {
		if strings.Contains(ownersStr, sep) {
			owners = strings.Split(ownersStr, sep)
			// Trim any whitespace
			for i, owner := range owners {
				owners[i] = strings.TrimSpace(owner)
			}
			return owners
		}
	}
	
	// If no separator found, treat as single owner
	return []string{ownersStr}
}

// ParseCustomFields processes custom fields string
func ParseCustomFields(customFieldsStr string) map[string]string {
	result := make(map[string]string)
	if customFieldsStr == "" {
		return result
	}
	
	// Split by semicolon for multiple fields
	fields := strings.Split(customFieldsStr, ";")
	for _, field := range fields {
		// Split by equals sign for key-value pairs
		if strings.Contains(field, "=") {
			parts := strings.SplitN(field, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				result[key] = value
			}
		}
	}
	
	return result
}