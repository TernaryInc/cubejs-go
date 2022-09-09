package cube_test

import (
	"encoding/json"
	"testing"

	cube "github.com/TernaryInc/cubejs-go"
	"github.com/stretchr/testify/assert"
)

func Test_DateRangeMarshalJSON(t *testing.T) {
	var battery = []struct {
		dateRange cube.DateRange
		expected  *string
	}{
		{
			cube.DateRange{RelativeRange: boxString("two weeks ago")},
			boxString(`"two weeks ago"`),
		},
		{
			cube.DateRange{AbsoluteRange: []string{"2021-04-20", "2021-04-21"}},
			boxString(`["2021-04-20","2021-04-21"]`),
		},
		{
			cube.DateRange{
				RelativeRange: boxString("two weeks ago"),
				AbsoluteRange: []string{"2021-04-20", "2021-04-21"},
			},
			nil,
		},
		{
			cube.DateRange{},
			nil,
		},
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

func Test_OrderTupleMarshalJSON(t *testing.T) {
	var battery = []struct {
		orderArray []cube.OrderTuple
		expected   string
	}{
		{
			[]cube.OrderTuple{{Key: "key", Order: cube.Order_Asc}},
			`[["key","asc"]]`,
		},
		{
			[]cube.OrderTuple{{Key: "key1", Order: cube.Order_Asc}, {Key: "key2", Order: cube.Order_Desc}},
			`[["key1","asc"],["key2","desc"]]`,
		},
		{
			nil,
			`null`,
		},
		{
			[]cube.OrderTuple{},
			`[]`,
		},
	}

	for _, tcase := range battery {
		var actual, err = json.Marshal(tcase.orderArray)
		assert.Nil(t, err)
		assert.Equal(t, tcase.expected, string(actual))
	}
}

func Test_OrderTupleUnmarshalJSON(t *testing.T) {
	var battery = []struct {
		jsonstr  string
		expected []cube.OrderTuple
	}{
		{
			`[["key","asc"]]`,
			[]cube.OrderTuple{{Key: "key", Order: cube.Order_Asc}},
		},
		{
			`[["key1","asc"],["key2","desc"]]`,
			[]cube.OrderTuple{{Key: "key1", Order: cube.Order_Asc}, {Key: "key2", Order: cube.Order_Desc}},
		},
		{
			`null`,
			nil,
		},
		{
			`[]`,
			[]cube.OrderTuple{},
		},
	}

	for _, tcase := range battery {
		var actual []cube.OrderTuple
		var err = json.Unmarshal([]byte(tcase.jsonstr), &actual)
		assert.Nil(t, err)
		assert.Equal(t, tcase.expected, actual)
	}
}

func boxString(x string) *string { return &x }
