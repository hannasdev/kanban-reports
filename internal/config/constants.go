package config

const (
	// MaxSuggestionsDisplay limits the number of CSV file suggestions shown to users
	MaxSuggestionsDisplay = 3
	
	// DefaultPeriodType is the default time period for metrics grouping
	DefaultPeriodType = "month"
	
	// DefaultAdHocFilter is the default ad-hoc request filtering behavior
	DefaultAdHocFilter = "include"
	
	// DefaultFilterField is the default date field used for filtering
	DefaultFilterField = "completed_at"
	
	// DefaultDelimiter is the default CSV delimiter setting
	DefaultDelimiter = "auto"
	
	// DateFormat is the expected date format for command-line date inputs
	DateFormat = "2006-01-02"
	
	// HoursPerDay is used for end-of-day calculations
	HoursPerDay = 23
	MinutesPerHour = 59
	SecondsPerMinute = 59
)