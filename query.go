package cube

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	cubeLoadPath = "/cubejs-api/v1/load"

	// Maximum duration a query should be retried for
	maximumQueryDuration = time.Duration(time.Minute * 30)

	Order_Asc  Order = "asc"
	Order_Desc Order = "desc"

	// Boolean logical operators currently unsupported: https://cube.dev/docs/query-format#filters-operators
	Operator_Equals               Operator = "equals"
	Operator_NotEquals            Operator = "notEquals"
	Operator_Contains             Operator = "contains"
	Operator_NotContains          Operator = "notContains"
	Operator_GreaterThan          Operator = "gt"
	Operator_GreaterThanOrEqualTo Operator = "gte"
	Operator_LessThan             Operator = "lt"
	Operator_LessThanOrEqualTo    Operator = "lte"
	Operator_Set                  Operator = "set"
	Operator_NotSet               Operator = "notSet"
	Operator_InDateRange          Operator = "inDateRange"
	Operator_NotInDateRange       Operator = "notInDateRange"
	Operator_BeforeDate           Operator = "beforeDate"
	Operator_AfterDate            Operator = "afterDate"

	Granularity_Second  Granularity = "second"
	Granularity_Minute  Granularity = "minute"
	Granularity_Hour    Granularity = "hour"
	Granularity_Day     Granularity = "day"
	Granularity_Week    Granularity = "week"
	Granularity_Month   Granularity = "month"
	Granularity_Quarter Granularity = "quarter"
	Granularity_Year    Granularity = "year"
)

// https://cube.dev/docs/@cubejs-client-core#types-time-dimension-granularity
type Granularity string

// https://cube.dev/docs/@cubejs-client-core#types-filter-operator
type Operator string

// https://cube.dev/docs/@cubejs-client-core#order
type Order string

type OrderTuple struct {
	Key   string
	Order Order
}

func (t OrderTuple) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string{t.Key, string(t.Order)})
}

func (t *OrderTuple) UnmarshalJSON(byts []byte) error {
	tmp := [2]string{}

	err := json.Unmarshal(byts, &tmp)
	if err != nil {
		return fmt.Errorf("unmarshal bytes into string array: %w", err)
	}

	t.Key = tmp[0]
	t.Order = Order(tmp[1])
	return nil
}

type requestBody struct {
	Query Query `json:"query"`
}

// Query represents a query that can be issued to a Cube server via the client.
type Query struct {
	Measures       []string        `json:"measures,omitempty"`
	TimeDimensions []TimeDimension `json:"timeDimensions,omitempty"`
	Order          []OrderTuple    `json:"order,omitempty"`
	Limit          int             `json:"limit,omitempty"`
	Filters        []Filter        `json:"filters,omitempty"`
	Dimensions     []string        `json:"dimensions,omitempty"`
	Ungrouped      bool            `json:"ungrouped"`
}

// https://cube.dev/docs/query-format#time-dimensions-format
type TimeDimension struct {
	Dimension   string      `json:"dimension"`
	DateRange   DateRange   `json:"dateRange"`
	Granularity Granularity `json:"granularity,omitempty"`
}

// https://cube.dev/docs/@cubejs-client-core#order
type Filter struct {
	Member   string   `json:"member"`
	Operator Operator `json:"operator"`
	Values   []string `json:"values"`
}

// ResponseMetadata returns metadata that appears in the response from the
// Cube API that is not the requested data.
type ResponseMetadata struct {
	Query      interface{} `json:"query"`
	Annotation interface{} `json:"annotation"`
}

type responseBody struct {
	Data  json.RawMessage `json:"data"`
	Error string          `json:"error"`
	ResponseMetadata
}

// DateRange represents the (string|[]string) date range type in the Cube query format.
// https://cube.dev/docs/query-format
// https://cube.dev/docs/@cubejs-client-core#date-range
//
// This is a union type and only one field should be set.
type DateRange struct {
	RelativeRange *string
	AbsoluteRange []string
}

// RelativeDateRange returns a DateRange with the RelativeRange field set to the input string
// Example arguments: "last 7 days", "this month", "1 hour ago"
func RelativeDateRange(dateRange string) DateRange {
	return DateRange{
		RelativeRange: &dateRange,
	}
}

// MarshalJSON marshals the input DateRange object; only one of the fields (i.e. RelativeRange, AbsoluteRange) will be marshalled as a top-level JSON value, depending on which is set.
func (d DateRange) MarshalJSON() ([]byte, error) {
	if d.RelativeRange != nil && d.AbsoluteRange == nil {
		return json.Marshal(d.RelativeRange)
	} else if len(d.AbsoluteRange) > 0 && d.RelativeRange == nil {
		return json.Marshal(d.AbsoluteRange)
	} else {
		return []byte{}, fmt.Errorf("invalid date range: exactly one field must be set: %+v", d)
	}
}
