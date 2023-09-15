package cube_test

import (
	"net/url"
	"testing"

	cube "github.com/TernaryInc/cubejs-go"
	"github.com/stretchr/testify/assert"
)

func Test_CubeRESTAPIBasePath(t *testing.T) {
	var battery = []struct {
		cubeURL  *url.URL
		expected string
	}{
		{
			&url.URL{Host: "http://localhost:4000", Path: ""},
			`/cubejs-api/v1`,
		},
		{
			&url.URL{Host: "my.cube", Path: ""},
			`/cubejs-api/v1`,
		},
		{
			&url.URL{Host: "my.cube", Path: "/analytics"},
			`/analytics/cubejs-api/v1`,
		},
	}

	for _, tcase := range battery {
		client := cube.NewClient(*tcase.cubeURL, nil)
		assert.Equal(t, client.CubeRESTAPIBasePath(), tcase.expected)
	}
}
