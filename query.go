package cube

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"go.uber.org/ratelimit"
	"go.uber.org/zap"

	"gitlab.com/ternary-app/backend/pkg/telemetry"
	"gitlab.com/ternary-app/backend/pkg/util/logging"
)

const (
	ASC  = "asc"
	DESC = "desc"

	DAY   = "day"
	MONTH = "month"
	WEEK  = "week"

	GCP_COST         = "GCPBillingLabels.cost"
	GCP_CREDITS      = "GCPBillingLabels.credits"
	GCP_SKU_CATEGORY = "GCPBillingLabels.category"
	GCP_TIMESTAMP    = "GCPBillingLabels.timestamp"

	EQUALS = "equals"

	cubeLoadPath = "/cubejs-api/v1/load"

	maximumQueryDuration = time.Duration(time.Minute * 30)

	granularity_Daily   Granularity = "daily"
	granularity_Weekly  Granularity = "weekly"
	granularity_Monthly Granularity = "monthly"
)

type loadBody struct {
	Query CubeQuery `json:"query"`
}

type CubeQuery struct {
	Measures       []string          `json:"measures,omitempty"`
	TimeDimensions []TimeDimension   `json:"timeDimensions,omitempty"`
	Order          map[string]string `json:"order,omitempty"`
	Limit          int               `json:"limit,omitempty"`
	Filters        []Filter          `json:"filters,omitempty"`
	Dimensions     []string          `json:"dimensions,omitempty"`
}

type TimeDimension struct {
	Dimension   string      `json:"dimension"`
	DateRange   interface{} `json:"dateRange"`
	Granularity string      `json:"granularity"`
}

type Filter struct {
	Member   string   `json:"member"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

type LoadResponse struct {
	Data  []LoadData `json:"data"`
	Error string     `json:"error"`
}

type LoadData map[string]interface{}

type CubeError struct {
	ErrorMessage string
	StatusCode   int
}

type Granularity string

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

func (c *CubeClient) Load(ctx context.Context, query CubeQuery) (*LoadResponse, error) {
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("invalid Cube query: %w", err)
	}

	var beginTime = time.Now()
	loadBody := loadBody{query}

	token, err := c.getToken()
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}

	marshaledLoadBody, err := json.Marshal(loadBody)
	if err != nil {
		return nil, fmt.Errorf("marshal load body: %w", err)
	}

	var (
		attempt            = 0
		continueWaitString = "Continue wait"
		// Do not spam Cube server with requests that will likely take some time
		// TODO(bruce): Replace with exponential backoff rate limiter
		limiter = ratelimit.New(1, ratelimit.Per(time.Minute))
	)

	for true {
		var loadResponse LoadResponse
		var response *http.Response
		attempt++

		req, err := http.NewRequestWithContext(ctx, "POST", c.cubeURL+cubeLoadPath, bytes.NewBuffer(marshaledLoadBody))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)

		limiter.Take()
		response, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("do request: %w", err)
		}

		defer func(body io.ReadCloser) {
			_ = body.Close()
		}(response.Body)

		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("read response bytes: %w", err)
		}

		if response.StatusCode >= 400 {
			return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, strings.TrimSpace(string(bodyBytes)))
		}

		bodyReader := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		err = json.NewDecoder(bodyReader).Decode(&loadResponse)
		if err != nil {
			preview := strings.TrimSpace(string(bodyBytes))
			rawLen := len(preview)
			if rawLen > 1024 {
				preview = preview[:1024]
			}

			return nil, fmt.Errorf("decode response json (%s...+%d more): %w", string(preview), rawLen-len(preview), err)
		}

		currentTime := time.Now()

		if loadResponse.Error == "" {
			telemetry.FromCtx(ctx).Gauge("backend.cube.query_duration_seconds", currentTime.Sub(beginTime).Seconds(), nil)
			return &loadResponse, nil
		} else if loadResponse.Error != continueWaitString {
			return nil, &CubeError{ErrorMessage: loadResponse.Error, StatusCode: response.StatusCode}
		} else if currentTime.Sub(beginTime) > maximumQueryDuration {
			return nil, &CubeError{ErrorMessage: fmt.Sprintf("Maximum query duration exceeded.  Number of attempts: %d.  Duration: %+v", attempt, currentTime.Sub(beginTime)), StatusCode: http.StatusInternalServerError}
		}

		logging.FromCtx(ctx).Info("Continue wait -- automatically retrying query", zap.Int("attempt", attempt))
	}

	return nil, &CubeError{ErrorMessage: "Unable to load Cube query", StatusCode: http.StatusInternalServerError}
}

func MonthlySpend(dateRange string, dimensions []string, filters []Filter) CubeQuery {
	return PeriodicSpend(dateRange, dimensions, MONTH, filters)
}

func DailySpend(dateRange string, dimensions []string, filters []Filter) CubeQuery {
	return PeriodicSpend(dateRange, dimensions, DAY, filters)
}

func WeeklySpend(dateRange string, dimensions []string, filters []Filter) CubeQuery {
	return PeriodicSpend(dateRange, dimensions, WEEK, filters)
}

func PeriodicSpend(dateRange string, dimensions []string, granularity string, filters []Filter) CubeQuery {
	measures := []string{
		GCP_COST,
		GCP_CREDITS,
	}
	timeDimensions := []TimeDimension{
		{
			Dimension:   GCP_TIMESTAMP,
			Granularity: granularity,
			DateRange:   dateRange,
		},
	}
	order := map[string]string{GCP_COST: DESC}

	return CubeQuery{
		Measures:       measures,
		TimeDimensions: timeDimensions,
		Order:          order,
		Limit:          10_000,
		Filters:        filters,
		Dimensions:     dimensions,
	}
}

func CategorySpend(dateRange string, granularity Granularity, dimension []string, filter []Filter) CubeQuery {
	switch granularity {
	case granularity_Daily:
		return DailySpend(dateRange, dimension, filter)
	case granularity_Weekly:
		return WeeklySpend(dateRange, dimension, filter)
	case granularity_Monthly:
		return MonthlySpend(dateRange, dimension, filter)
	default:
		panic(fmt.Sprintf("unknown granularity: %s", granularity))
	}
}

func FormatDateRange(begin, end time.Time) string {
	return fmt.Sprintf("%v to %v", begin.Format("2006-01-02"), end.Format("2006-01-02"))
}
