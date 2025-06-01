package menu

import (
	"fmt"
	"strings"
)

// QuitError represents a user-initiated quit
type QuitError struct {
	Message string
}

func (e QuitError) Error() string {
	return e.Message
}

// IsQuitCommand checks if the input is a quit command
func IsQuitCommand(input string) bool {
	trimmed := strings.ToLower(strings.TrimSpace(input))
	quitCommands := []string{"q", "quit", "exit", "bye"}
	
	for _, cmd := range quitCommands {
		if trimmed == cmd {
			return true
		}
	}
	return false
}

// HandleQuit returns a QuitError if the input is a quit command
func HandleQuit(input string) error {
	if IsQuitCommand(input) {
		return QuitError{Message: "User requested to quit"}
	}
	return nil
}

// ShowQuitHelp displays quit help information
func ShowQuitHelp() {
	fmt.Println("\n💡 Tip: Type 'q', 'quit', 'exit', or 'bye' at any time to exit")
}