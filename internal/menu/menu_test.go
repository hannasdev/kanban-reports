package menu

import (
	"strings"
	"testing"

	"github.com/hannasdev/kanban-reports/internal/config"
	"github.com/hannasdev/kanban-reports/internal/reports"
)

// TestInput simulates user input for testing
func createTestMenu(input string) *Menu {
	reader := strings.NewReader(input)
	writer := &strings.Builder{}
	return NewMenuWithIO(reader, writer)
}

func TestQuitCommands(t *testing.T) {
	quitCommands := []string{"q", "quit", "exit", "bye", "Q", "QUIT", "Exit", "BYE"}
	
	for _, cmd := range quitCommands {
		t.Run("Quit_with_"+cmd, func(t *testing.T) {
			menu := createTestMenu(cmd + "\n")
			
			_, err := menu.readInput("Test prompt: ")
			
			// Should return a QuitError
			if err == nil {
				t.Errorf("Expected QuitError for command '%s', got nil", cmd)
			}
			
			quitErr, ok := err.(QuitError)
			if !ok {
				t.Errorf("Expected QuitError for command '%s', got %T: %v", cmd, err, err)
			}
			
			if quitErr.Message != "User requested to quit" {
				t.Errorf("Expected quit message, got: %s", quitErr.Message)
			}
		})
	}
}

func TestChooseMode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantMode bool // true for metrics, false for reports
		wantErr  bool
	}{
		{"Select reports mode", "1\n", false, false},
		{"Select metrics mode", "2\n", true, false},
		{"Invalid input then valid", "3\n1\n", false, false},
		{"Quit command", "q\n", false, true},
		{"Empty input then valid", "\n1\n", false, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu := createTestMenu(tt.input)
			
			isMetrics, err := menu.chooseMode()
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for input '%s', got nil", tt.input)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Expected no error for input '%s', got: %v", tt.input, err)
			}
			
			if isMetrics != tt.wantMode {
				t.Errorf("Expected mode %v for input '%s', got %v", tt.wantMode, tt.input, isMetrics)
			}
		})
	}
}

func TestConfigureReports(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType reports.ReportType
		wantErr  bool
	}{
		{"Select contributor", "1\n", reports.ReportTypeContributor, false},
		{"Select epic", "2\n", reports.ReportTypeEpic, false},
		{"Select product area", "3\n", reports.ReportTypeProductArea, false},
		{"Select team", "4\n", reports.ReportTypeTeam, false},
		{"Invalid then valid", "5\n1\n", reports.ReportTypeContributor, false},
		{"Quit command", "quit\n", "", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu := createTestMenu(tt.input)
			cfg := &config.Config{}
			
			err := menu.configureReports(cfg)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for input '%s', got nil", tt.input)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Expected no error for input '%s', got: %v", tt.input, err)
			}
			
			if cfg.ReportType != tt.wantType {
				t.Errorf("Expected report type %v, got %v", tt.wantType, cfg.ReportType)
			}
		})
	}
}

func TestConfigureLastNDays(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantDays int
		wantErr  bool
	}{
		{"Valid 7 days", "7\n", 7, false},
		{"Valid 30 days", "30\n", 30, false},
		{"Invalid then valid", "abc\n7\n", 7, false},
		{"Negative then valid", "-5\n7\n", 7, false},
		{"Zero then valid", "0\n7\n", 7, false},
		{"Quit command", "bye\n", 0, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu := createTestMenu(tt.input)
			cfg := &config.Config{}
			
			err := menu.configureLastNDays(cfg)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for input '%s', got nil", tt.input)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Expected no error for input '%s', got: %v", tt.input, err)
			}
			
			if cfg.LastNDays != tt.wantDays {
				t.Errorf("Expected %d days, got %d", tt.wantDays, cfg.LastNDays)
			}
		})
	}
}

func TestConfigureSpecificRange(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantStart string
		wantEnd   string
		wantErr   bool
	}{
		{
			name:      "Valid date range",
			input:     "2024-01-01\n2024-01-31\n",
			wantStart: "2024-01-01",
			wantEnd:   "2024-01-31",
			wantErr:   false,
		},
		{
			name:      "Invalid start date then valid",
			input:     "invalid\n2024-01-01\n2024-01-31\n",
			wantStart: "2024-01-01",
			wantEnd:   "2024-01-31",
			wantErr:   false,
		},
		{
			name:    "Quit on start date",
			input:   "q\n",
			wantErr: true,
		},
		{
			name:    "Quit on end date",
			input:   "2024-01-01\nquit\n",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu := createTestMenu(tt.input)
			cfg := &config.Config{}
			
			err := menu.configureSpecificRange(cfg)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for input '%s', got nil", tt.input)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Expected no error for input '%s', got: %v", tt.input, err)
				return
			}
			
			startStr := cfg.StartDate.Format("2006-01-02")
			endStr := cfg.EndDate.Format("2006-01-02")
			
			if startStr != tt.wantStart {
				t.Errorf("Expected start date %s, got %s", tt.wantStart, startStr)
			}
			
			if endStr != tt.wantEnd {
				t.Errorf("Expected end date %s, got %s", tt.wantEnd, endStr)
			}
		})
	}
}

func TestIsQuitCommand(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"q", true},
		{"quit", true},
		{"exit", true},
		{"bye", true},
		{"Q", true},
		{"QUIT", true},
		{"Exit", true},
		{" q ", true},      // with spaces
		{"\tquit\n", true}, // with whitespace
		{"query", false},   // contains but not exact
		{"1", false},
		{"", false},
		{"help", false},
	}
	
	for _, tt := range tests {
		t.Run("Input_"+tt.input, func(t *testing.T) {
			if got := IsQuitCommand(tt.input); got != tt.want {
				t.Errorf("IsQuitCommand(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestQuitHandling_Integration(t *testing.T) {
	helper := NewTestHelper()
	defer helper.Cleanup()
	
	// Create a temp CSV file for tests that need it
	tmpFile := helper.CreateTempCSV(t, "")
	
	tests := []struct {
		name  string
		input string
	}{
		{"Quit at file selection", "q\n"},
		{"Quit at mode selection", tmpFile + "\nquit\n"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu := createTestMenu(tt.input)
			
			_, err := menu.Run()
			
			if err == nil {
				t.Errorf("Expected QuitError for '%s', got nil", tt.name)
			}
			
			if _, ok := err.(QuitError); !ok {
				t.Errorf("Expected QuitError for '%s', got %T: %v", tt.name, err, err)
			}
		})
	}
}