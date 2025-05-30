package reports

import (
	"fmt"
)

// ReportType defines the type of report to generate
type ReportType string

const (
	// ReportTypeContributor generates report by contributor
	ReportTypeContributor ReportType = "contributor"
	// ReportTypeEpic generates report by epic
	ReportTypeEpic ReportType = "epic"
	// ReportTypeProductArea generates report by product area
	ReportTypeProductArea ReportType = "product-area"
	// ReportTypeTeam generates report by team
	ReportTypeTeam ReportType = "team"
)

// Validation function for ReportType
func (rt ReportType) IsValid() bool {
	switch rt {
	case ReportTypeContributor, ReportTypeEpic, ReportTypeProductArea, ReportTypeTeam:
		return true
	}
	return false
}

// Parse strings into ReportType
func ParseReportType(s string) (ReportType, error) {
	rt := ReportType(s)
	if !rt.IsValid() {
		return "", fmt.Errorf("invalid report type: %s", s)
	}
	return rt, nil
}
