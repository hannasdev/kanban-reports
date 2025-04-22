package metrics

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

// PeriodType defines the time period for grouping metrics
type PeriodType string

const (
    // PeriodTypeWeek groups metrics by week
    PeriodTypeWeek PeriodType = "week"
    // PeriodTypeMonth groups metrics by month
    PeriodTypeMonth PeriodType = "month"
)