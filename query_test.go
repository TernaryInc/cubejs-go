package cube_test

import (
	"testing"

	cube "github.com/TernaryInc/cubejs-go"
	"github.com/stretchr/testify/assert"
)

// var (
// 	dailyRange   = "1 day ago"
// 	weeklyRange  = "this week"
// 	monthlyRange = "this month"
// )

// func Test_PeriodicSpend(t *testing.T) {
// 	tcases := []struct {
// 		dateRange   string
// 		dimensions  []string
// 		granularity string
// 		filters     []Filter
// 	}{
// 		{
// 			dateRange:   weeklyRange,
// 			dimensions:  []string{GCP_SKU_CATEGORY},
// 			granularity: WEEK,
// 		},
// 		{
// 			dateRange:   monthlyRange,
// 			dimensions:  []string{GCP_SKU_CATEGORY},
// 			granularity: MONTH,
// 		},
// 		{
// 			dateRange:   dailyRange,
// 			granularity: DAY,
// 		},
// 	}

// 	for index, tcase := range tcases {
// 		periodicSpend := PeriodicSpend(tcase.dateRange, tcase.dimensions, tcase.granularity, tcase.filters)

// 		assert.Equalf(t, periodicSpend.Dimensions, tcase.dimensions, "loop: %v - dimension check", index)
// 		assert.Equalf(t, periodicSpend.TimeDimensions[0].DateRange, tcase.dateRange, "loop: %v - datarange check", index)
// 		assert.Equalf(t, periodicSpend.TimeDimensions[0].Granularity, tcase.granularity, "loop: %v - granularity check", index)
// 		assert.Equalf(t, periodicSpend.Filters, tcase.filters, "loop: %v - filter check", index)
// 	}

// }

// func Test_Weekly_Spend(t *testing.T) {
// 	tcases := []struct {
// 		dateRange  string
// 		dimensions []string
// 		filters    []Filter
// 	}{
// 		{
// 			dateRange:  weeklyRange,
// 			dimensions: []string{GCP_SKU_CATEGORY},
// 		},
// 		{
// 			dateRange:  monthlyRange,
// 			dimensions: []string{GCP_SKU_CATEGORY},
// 		},
// 		{
// 			dateRange: dailyRange,
// 		},
// 	}

// 	for index, tcase := range tcases {
// 		periodicSpend := WeeklySpend(tcase.dateRange, tcase.dimensions, tcase.filters)

// 		assert.Equalf(t, periodicSpend.Dimensions, tcase.dimensions, "loop: %v - dimension check", index)
// 		assert.Equalf(t, periodicSpend.TimeDimensions[0].DateRange, tcase.dateRange, "loop: %v - datarange check", index)
// 		assert.Equalf(t, periodicSpend.TimeDimensions[0].Granularity, WEEK, "loop: %v - granularity check", index)
// 		assert.Equalf(t, periodicSpend.Filters, tcase.filters, "loop: %v - filter check", index)
// 	}
// }

// func Test_Monthly_Spend(t *testing.T) {
// 	tcases := []struct {
// 		dateRange  string
// 		dimensions []string
// 		filters    []Filter
// 	}{
// 		{
// 			dateRange:  weeklyRange,
// 			dimensions: []string{GCP_SKU_CATEGORY},
// 		},
// 		{
// 			dateRange:  monthlyRange,
// 			dimensions: []string{GCP_SKU_CATEGORY},
// 		},
// 		{
// 			dateRange: dailyRange,
// 		},
// 	}

// 	for index, tcase := range tcases {
// 		periodicSpend := MonthlySpend(tcase.dateRange, tcase.dimensions, tcase.filters)

// 		assert.Equalf(t, periodicSpend.Dimensions, tcase.dimensions, "loop: %v - dimension check", index)
// 		assert.Equalf(t, periodicSpend.TimeDimensions[0].DateRange, tcase.dateRange, "loop: %v - datarange check", index)
// 		assert.Equalf(t, periodicSpend.TimeDimensions[0].Granularity, MONTH, "loop: %v - granularity check", index)
// 		assert.Equalf(t, periodicSpend.Filters, tcase.filters, "loop: %v - filter check", index)
// 	}
// }

func boxString(x string) *string { return &x }

func Test_DateRangeMarshalJSON(t *testing.T) {
	var battery = []struct {
		dateRange cube.DateRange
		expected  *string
	}{
		{cube.DateRange{RelativeRange: boxString("two weeks ago")}, boxString(`"two weeks ago"`)},
		{cube.DateRange{AbsoluteRange: []string{"2021-04-20", "2021-04-21"}}, boxString(`["2021-04-20","2021-04-21"]`)},
		{cube.DateRange{RelativeRange: boxString("two weeks ago"), AbsoluteRange: []string{"2021-04-20", "2021-04-21"}}, nil},
		{cube.DateRange{}, nil},
	}

	for _, tcase := range battery {
		var actual, err = tcase.dateRange.MarshalJSON()
		if tcase.expected == nil {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, *tcase.expected, string(actual))
		}
	}
}

// func Test_ValidateTimeDimension(t *testing.T) {
// 	var battery = []struct {
// 		dateRange interface{}
// 		valid     bool
// 	}{
// 		{
// 			dateRange: []string{"length one"},
// 			valid:     false,
// 		},
// 		{
// 			dateRange: []string{"length", "is", "three"},
// 			valid:     false,
// 		},
// 		{
// 			dateRange: []string{"length", "two"},
// 			valid:     true,
// 		},
// 		{
// 			dateRange: "just a string",
// 			valid:     true,
// 		},
// 		{
// 			dateRange: 1337,
// 			valid:     false,
// 		},
// 		{
// 			dateRange: []float64{1, 2},
// 			valid:     false,
// 		},
// 		{
// 			dateRange: []interface{}{"asdf", 2},
// 			valid:     false,
// 		},
// 		{
// 			dateRange: []interface{}{"also", "two"},
// 			valid:     true,
// 		},
// 	}

// 	for _, tcase := range battery {
// 		var timeDimension = cube.TimeDimension{
// 			DateRange: tcase.dateRange,
// 		}

// 		var err = timeDimension.Validate()

// 		if tcase.valid {
// 			assert.Nil(t, err)
// 		} else {
// 			assert.Error(t, err)
// 		}
// 	}
// }

// func Test_Daily_Spend(t *testing.T) {
// 	tcases := []struct {
// 		dateRange  string
// 		dimensions []string
// 		filters    []Filter
// 	}{
// 		{
// 			dateRange:  weeklyRange,
// 			dimensions: []string{GCP_SKU_CATEGORY},
// 		},
// 		{
// 			dateRange:  monthlyRange,
// 			dimensions: []string{GCP_SKU_CATEGORY},
// 		},
// 		{
// 			dateRange: dailyRange,
// 		},
// 	}

// 	for index, tcase := range tcases {
// 		periodicSpend := DailySpend(tcase.dateRange, tcase.dimensions, tcase.filters)

// 		assert.Equalf(t, periodicSpend.Dimensions, tcase.dimensions, "loop: %v - dimension check", index)
// 		assert.Equalf(t, periodicSpend.TimeDimensions[0].DateRange, tcase.dateRange, "loop: %v - datarange check", index)
// 		assert.Equalf(t, periodicSpend.TimeDimensions[0].Granularity, DAY, "loop: %v - granularity check", index)
// 		assert.Equalf(t, periodicSpend.Filters, tcase.filters, "loop: %v - filter check", index)
// 	}
// }

// func Test_CategorySpend(t *testing.T) {
// 	tcases := []struct {
// 		dateRange           string
// 		dimensions          []string
// 		filters             []Filter
// 		granularity         Granularity
// 		espectedGranularity string
// 	}{
// 		{
// 			dateRange:           dailyRange,
// 			dimensions:          []string{GCP_SKU_CATEGORY},
// 			granularity:         "daily",
// 			espectedGranularity: DAY,
// 		},
// 		{
// 			dateRange:           weeklyRange,
// 			dimensions:          []string{GCP_SKU_CATEGORY},
// 			granularity:         "weekly",
// 			espectedGranularity: WEEK,
// 		},
// 		{
// 			dateRange:           monthlyRange,
// 			granularity:         "monthly",
// 			espectedGranularity: MONTH,
// 		},
// 	}

// 	for index, tcase := range tcases {
// 		cSpend := CategorySpend(tcase.dateRange, tcase.granularity, tcase.dimensions, tcase.filters)
// 		assert.Equalf(t, cSpend.Dimensions, tcase.dimensions, "loop: %v - dimension check", index)
// 		assert.Equalf(t, cSpend.TimeDimensions[0].DateRange, tcase.dateRange, "loop: %v - datarange check", index)
// 		assert.Equalf(t, cSpend.TimeDimensions[0].Granularity, tcase.espectedGranularity, "loop: %v - granularity check", index)
// 		assert.Equalf(t, cSpend.Filters, tcase.filters, "loop: %v - filter check", index)
// 	}
// }

// func Test_FormatDateRange(t *testing.T) {
// 	tcases := []struct {
// 		FromResult     time.Time
// 		ToResult       time.Time
// 		ExpectedFormat string
// 	}{
// 		{
// 			FromResult:     time.Date(2021, 11, 8, 0, 0, 0, 0, time.UTC),
// 			ToResult:       time.Date(2021, 11, 14, 0, 0, 0, 0, time.UTC),
// 			ExpectedFormat: "2021-11-08 to 2021-11-14",
// 		},
// 		{
// 			FromResult:     time.Date(2021, 10, 18, 0, 0, 0, 0, time.UTC),
// 			ToResult:       time.Date(2021, 10, 24, 0, 0, 0, 0, time.UTC),
// 			ExpectedFormat: "2021-10-18 to 2021-10-24",
// 		},
// 		{
// 			FromResult:     time.Date(2021, 02, 15, 0, 0, 0, 0, time.UTC),
// 			ToResult:       time.Date(2021, 02, 21, 0, 0, 0, 0, time.UTC),
// 			ExpectedFormat: "2021-02-15 to 2021-02-21",
// 		},
// 		{
// 			FromResult:     time.Date(2021, 02, 22, 0, 0, 0, 0, time.UTC),
// 			ToResult:       time.Date(2021, 02, 28, 0, 0, 0, 0, time.UTC),
// 			ExpectedFormat: "2021-02-22 to 2021-02-28",
// 		},
// 	}

// 	for index, tcase := range tcases {
// 		actualFormat := FormatDateRange(tcase.FromResult, tcase.ToResult)
// 		assert.Equalf(t, tcase.ExpectedFormat, actualFormat, "%v loop", index)
// 	}
// }
