package cube

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	Order_Asc  Order = "asc"
	Order_Desc Order = "desc"

	// TODO(bruce): Test unary operators?
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
	// Boolean logical operators currently unsupported: https://cube.dev/docs/query-format#filters-operators

	cubeLoadPath = "/cubejs-api/v1/load"

	// Maximum duration a query should be retried for
	maximumQueryDuration = time.Duration(time.Minute * 30)

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

type requestBody struct {
	Query CubeQuery `json:"query"`
}

// CubeQuery represents a query that can be issued to a Cube server via the client.
type CubeQuery struct {
	Measures       []string        `json:"measures,omitempty"`
	TimeDimensions []TimeDimension `json:"timeDimensions,omitempty"`
	// TODO: Why is this a map[string]string?
	Order      map[string]string `json:"order,omitempty"`
	Limit      int               `json:"limit,omitempty"`
	Filters    []Filter          `json:"filters,omitempty"`
	Dimensions []string          `json:"dimensions,omitempty"`
}

// https://cube.dev/docs/query-format#time-dimensions-format
type TimeDimension struct {
	Dimension string `json:"dimension"`
	// TODO: Document interface{} or choose something else
	DateRange   interface{} `json:"dateRange"`
	Granularity Granularity `json:"granularity"`
}

// https://cube.dev/docs/@cubejs-client-core#order
type Filter struct {
	Member   string   `json:"member"`
	Operator Operator `json:"operator"`
	// TODO(Bruce): omitempty?
	Values []string `json:"values"`
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

type cubeError struct {
	ErrorMessage string
	StatusCode   int
}

// Validate determines whether the input query is valid.
func (query CubeQuery) Validate() error {
	for _, timeDimension := range query.TimeDimensions {
		if err := timeDimension.Validate(); err != nil {
			return fmt.Errorf("invalid time dimension: %w", err)
		}
	}

	return nil
}

// Validate determines whether the input time dimension is valid.
// The date range of a time dimension can either be a string or an array of strings of length two.
func (timeDimension TimeDimension) Validate() error {
	if _, ok := timeDimension.DateRange.(string); ok {
		return nil
	} else if arr, ok := timeDimension.DateRange.([]string); ok {
		if len(arr) != 2 {
			return fmt.Errorf("date range with array type must have length two, got %d", len(arr))
		}

		return nil
	} else if arr, ok := timeDimension.DateRange.([]interface{}); ok {
		if len(arr) != 2 {
			return fmt.Errorf("date range with array type must have length two")
		}

		for _, arrayElement := range arr {
			if _, ok := arrayElement.(string); !ok {
				return fmt.Errorf("date range with array type must have string entries.  got %+v, type %T", arrayElement, arrayElement)
			}
		}

		return nil
	}

	return fmt.Errorf("unsupported type for time dimension date range (value: %+v) (type: %T)", timeDimension.DateRange, timeDimension.DateRange)
}

func (ce *cubeError) Error() string {
	return fmt.Sprintf("StatusCode: %v ErrorMessage: %v", ce.StatusCode, ce.ErrorMessage)
}
