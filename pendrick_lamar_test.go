package cube_test

import (
	"context"
	"net/url"
	"testing"

	cube "github.com/TernaryInc/cubejs-go"
	"github.com/stretchr/testify/assert"
)

func TestPendrick(t *testing.T) {
	var u = url.URL{
		Scheme: "http",
		Host:   "localhost:4000",
	}

	var cubeClient = cube.NewCubeClient(u)

	var cubeQuery = cube.CubeQuery{
		Measures: []string{"GCPBillingDaily.cost"},
	}
	var results []int

	var _, err = cubeClient.Load(context.Background(), cubeQuery, &results)
	assert.Nil(t, err)
}
