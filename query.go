package cube

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	Order_Asc  Order = "asc"
	Order_Desc Order = "desc"

	// TODO: Pick a link
	// https://cube.dev/docs/query-format#filters-operators
	// https://cube.dev/docs/@cubejs-client-core#types-filter-operator
	// TODO: Test unary operators?
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

	// https://cube.dev/docs/@cubejs-client-core#types-time-dimension-granularity
	Granularity_Second  Granularity = "second"
	Granularity_Minute  Granularity = "minute"
	Granularity_Hour    Granularity = "hour"
	Granularity_Day     Granularity = "day"
	Granularity_Week    Granularity = "week"
	Granularity_Month   Granularity = "month"
	Granularity_Quarter Granularity = "quarter"
	Granularity_Year    Granularity = "year"
)

type Granularity string
type Operator string
type Order string

type requestBody struct {
	Query CubeQuery `json:"query"`
}

type CubeQuery struct {
	Measures       []string        `json:"measures,omitempty"`
	TimeDimensions []TimeDimension `json:"timeDimensions,omitempty"`
	// TODO: Why is this a map[string]string?
	Order      map[string]string `json:"order,omitempty"`
	Limit      int               `json:"limit,omitempty"`
	Filters    []Filter          `json:"filters,omitempty"`
	Dimensions []string          `json:"dimensions,omitempty"`
}

type TimeDimension struct {
	Dimension string `json:"dimension"`
	// TODO: Document interface{} or choose something else
	DateRange   interface{} `json:"dateRange"`
	Granularity string      `json:"granularity"`
}

type Filter struct {
	Member   string `json:"member"`
	Operator string `json:"operator"`
	// TODO(omitempty?)
	Values []string `json:"values"`
}

type ResponseMetadata struct {
	Query      interface{} `json:"query"`
	Annotation interface{} `json:"annotation"`
}

type ResponseBody struct {
	Data  json.RawMessage `json:"data"`
	Error string          `json:"error"`
	ResponseMetadata
}

// Is there anything we can do with struct tags and JSON (un?)marshaling
type LoadData map[string]interface{}

type CubeError struct {
	ErrorMessage string
	StatusCode   int
}

func (query CubeQuery) Validate() error {
	for _, timeDimension := range query.TimeDimensions {
		if err := timeDimension.Validate(); err != nil {
			return fmt.Errorf("invalid time dimension: %w", err)
		}
	}

	return nil
}

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

func (ce *CubeError) Error() string {
	return fmt.Sprintf("StatusCode: %v ErrorMessage: %v", ce.StatusCode, ce.ErrorMessage)
}

// func FormatDateRange(begin, end time.Time) string {
// 	return fmt.Sprintf("%v to %v", begin.Format("2006-01-02"), end.Format("2006-01-02"))
// }
