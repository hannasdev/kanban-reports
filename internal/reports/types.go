package reports

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