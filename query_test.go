package cube_test

import (
	"testing"

	cube "github.com/TernaryInc/cubejs-go"
	"github.com/stretchr/testify/assert"
)

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
