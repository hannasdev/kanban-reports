package metrics

import "fmt"

// MetricsType defines the type of metrics to generate
type MetricsType string

const (
    // MetricsTypeLeadTime generates lead time analysis by story point size
    MetricsTypeLeadTime MetricsType = "lead-time"
    // MetricsTypeThroughput generates throughput analysis over time
    MetricsTypeThroughput MetricsType = "throughput"
    // MetricsTypeFlow generates flow efficiency analysis
    MetricsTypeFlow MetricsType = "flow"
    // MetricsTypeEstimation generates estimation accuracy analysis
    MetricsTypeEstimation MetricsType = "estimation"
    // MetricsTypeAge generates current work item age analysis
    MetricsTypeAge MetricsType = "age"
    // MetricsTypeImprovement generates month-over-month improvement metrics
    MetricsTypeImprovement MetricsType = "improvement"
    // MetricsTypeAll generates all metrics reports
    MetricsTypeAll MetricsType = "all"
)

// Validate MetricsType
func (mt MetricsType) IsValid() bool {
    switch mt {
    case MetricsTypeLeadTime, MetricsTypeThroughput, MetricsTypeFlow, MetricsTypeEstimation, MetricsTypeAge, MetricsTypeImprovement, MetricsTypeAll:
        return true
    }
    return false
}

// Function to parse strings into MetricsType
func ParseMetricsType(s string) (MetricsType, error) {
    if s == "" {
        return "", nil // Empty is valid (no metrics)
    }
    mt := MetricsType(s)
    if !mt.IsValid() {
        return "", fmt.Errorf("invalid report type: %s", s)
    }
    return mt, nil
}

// PeriodType defines the time period for grouping metrics
type PeriodType string

const (
    // PeriodTypeWeek groups metrics by week
    PeriodTypeWeek PeriodType = "week"
    // PeriodTypeMonth groups metrics by month
    PeriodTypeMonth PeriodType = "month"
)

func (pt PeriodType) IsValid() bool {
    switch pt {
    case PeriodTypeWeek, PeriodTypeMonth:
        return true
    }
    return false
}

func ParsePeriodType(s string) (PeriodType, error) {
    pt := PeriodType(s)
    if !pt.IsValid() {
        return "", fmt.Errorf("invalid period type: %s (must be one of: week, month)", s)
    }
    return pt, nil
}
