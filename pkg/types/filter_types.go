package types

import "fmt"

// AdHocFilterType defines how to handle ad-hoc requests
type AdHocFilterType string

const (
	// AdHocFilterInclude includes all items (default)
	AdHocFilterInclude AdHocFilterType = "include"
	// AdHocFilterExclude excludes ad-hoc requests
	AdHocFilterExclude AdHocFilterType = "exclude"
	// AdHocFilterOnly shows only ad-hoc requests
	AdHocFilterOnly AdHocFilterType = "only"
)

// IsValid checks if an AdHocFilterType is valid
func (aft AdHocFilterType) IsValid() bool {
	switch aft {
	case AdHocFilterInclude, AdHocFilterExclude, AdHocFilterOnly:
		return true
	}
	return false
}

// ParseAdHocFilterType converts a string to an AdHocFilterType with validation
func ParseAdHocFilterType(s string) (AdHocFilterType, error) {
	aft := AdHocFilterType(s)
	if !aft.IsValid() {
		return "", fmt.Errorf("invalid ad-hoc filter type: %s (must be one of: include, exclude, only)", s)
	}
	return aft, nil
}