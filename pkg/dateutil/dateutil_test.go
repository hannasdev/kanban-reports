package dateutil

import (
	"testing"
	"time"
)

func TestGetStartOfPeriod(t *testing.T) {
	// Create a test date: Wednesday, May 15, 2024, 14:30:45
	testDate := time.Date(2024, 5, 15, 14, 30, 45, 0, time.UTC)

	tests := []struct {
		name       string
		date       time.Time
		periodType string
		expected   time.Time
	}{
		{
			name:       "Start of week for Wednesday",
			date:       testDate,
			periodType: "week",
			// Should return Sunday, May 12, 2024, 14:30:45 (3 days earlier)
			expected: time.Date(2024, 5, 12, 14, 30, 45, 0, time.UTC),
		},
		{
			name:       "Start of month for mid-month",
			date:       testDate,
			periodType: "month",
			// Should return May 1, 2024, 00:00:00
			expected: time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:       "Start of week for Sunday",
			date:       time.Date(2024, 5, 12, 10, 0, 0, 0, time.UTC), // Sunday
			periodType: "week",
			// Should return the same date (already Sunday)
			expected: time.Date(2024, 5, 12, 10, 0, 0, 0, time.UTC),
		},
		{
			name:       "Start of month for first day",
			date:       time.Date(2024, 5, 1, 15, 30, 0, 0, time.UTC), // May 1st
			periodType: "month",
			// Should return May 1, 2024, 00:00:00
			expected: time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:       "Start of week for Saturday",
			date:       time.Date(2024, 5, 18, 12, 0, 0, 0, time.UTC), // Saturday
			periodType: "week",
			// Should return Sunday, May 12, 2024, 12:00:00 (6 days earlier)
			expected: time.Date(2024, 5, 12, 12, 0, 0, 0, time.UTC),
		},
		{
			name:       "Invalid period type defaults to month",
			date:       testDate,
			periodType: "invalid",
			// Should default to month behavior
			expected: time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetStartOfPeriod(tt.date, tt.periodType)
			if !got.Equal(tt.expected) {
				t.Errorf("GetStartOfPeriod() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFormatPeriod(t *testing.T) {
	// Create test dates
	testDates := []time.Time{
		time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),  // January 15, 2024
		time.Date(2024, 12, 25, 15, 30, 0, 0, time.UTC), // December 25, 2024
		time.Date(2023, 6, 5, 0, 0, 0, 0, time.UTC),     // June 5, 2023
	}

	tests := []struct {
		name       string
		date       time.Time
		periodType string
		expected   string
	}{
		{
			name:       "Format week for January",
			date:       testDates[0],
			periodType: "week",
			expected:   "2024-W03", // January 15, 2024 is in week 3
		},
		{
			name:       "Format month for January",
			date:       testDates[0],
			periodType: "month",
			expected:   "2024-01",
		},
		{
			name:       "Format week for December",
			date:       testDates[1],
			periodType: "week",
			expected:   "2024-W52", // December 25, 2024 is in week 52
		},
		{
			name:       "Format month for December",
			date:       testDates[1],
			periodType: "month",
			expected:   "2024-12",
		},
		{
			name:       "Format month for June 2023",
			date:       testDates[2],
			periodType: "month",
			expected:   "2023-06",
		},
		{
			name:       "Invalid period type defaults to month",
			date:       testDates[0],
			periodType: "invalid",
			expected:   "2024-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatPeriod(tt.date, tt.periodType)
			if got != tt.expected {
				t.Errorf("FormatPeriod() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetStartOfPeriod_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		date       time.Time
		periodType string
		expected   time.Time
	}{
		{
			name:       "End of year week",
			date:       time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC), // Tuesday
			periodType: "week",
			// Should return Sunday, December 29, 2024
			expected: time.Date(2024, 12, 29, 23, 59, 59, 0, time.UTC),
		},
		{
			name:       "Beginning of year month",
			date:       time.Date(2024, 1, 1, 0, 0, 1, 0, time.UTC),
			periodType: "month",
			// Should return January 1, 2024, 00:00:00
			expected: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:       "Leap year February",
			date:       time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC), // Leap year
			periodType: "month",
			// Should return February 1, 2024, 00:00:00
			expected: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetStartOfPeriod(tt.date, tt.periodType)
			if !got.Equal(tt.expected) {
				t.Errorf("GetStartOfPeriod() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFormatPeriod_WeekEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{
			name:     "First week of year",
			date:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "2024-W01",
		},
		{
			name:     "Last week of year",
			date:     time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC), // Thursday in last week of 2024
			expected: "2024-W52",
		},
		{
			name:     "Week that spans years",
			date:     time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC), // Sunday
			expected: "2023-W52",
		},
		{
			name:     "Date in week belonging to next year",
			date:     time.Date(2024, 12, 30, 0, 0, 0, 0, time.UTC), // Monday in first week of 2025
			expected: "2025-W01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatPeriod(tt.date, "week")
			if got != tt.expected {
				t.Errorf("FormatPeriod() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetStartOfPeriod_Timezone(t *testing.T) {
	// Test with different timezones
	est, _ := time.LoadLocation("America/New_York")
	pst, _ := time.LoadLocation("America/Los_Angeles")

	// Same moment in time, different timezones
	utcTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	estTime := utcTime.In(est)
	pstTime := utcTime.In(pst)

	tests := []struct {
		name       string
		date       time.Time
		periodType string
	}{
		{"UTC timezone", utcTime, "month"},
		{"EST timezone", estTime, "month"},
		{"PST timezone", pstTime, "month"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetStartOfPeriod(tt.date, tt.periodType)
			// The result should preserve the original timezone
			if got.Location() != tt.date.Location() {
				t.Errorf("GetStartOfPeriod() timezone = %v, want %v", got.Location(), tt.date.Location())
			}
		})
	}
}