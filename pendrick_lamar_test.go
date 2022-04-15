package cube_test

import (
	"context"
	"fmt"
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

	var nower = cube.TimeNower{}
	tokenGenerator, err := cube.NewCubeTokenGenerator(cube.CubeAPISecret, nower)
	assert.Nil(t, err)

	token, err := tokenGenerator.GenerateToken("SyYlNf2WE0nN757SyYoX")
	assert.Nil(t, err)

	var cubeClient = cube.NewCubeClient(u, &token.Token)

	var cubeQuery = cube.CubeQuery{
		Measures:   []string{"GCPBillingDaily.cost"},
		Dimensions: []string{"GCPBillingDaily.projectId"},
	}

	type QueryResult struct {
		Cost      float64 `json:"GCPBillingDaily.cost"`
		ProjectID string  `json:"GCPBillingDaily.projectId"`
	}

	var results []QueryResult

	err = cubeClient.Load(context.Background(), cubeQuery, &results)
	assert.Nil(t, err)

	fmt.Printf("%+v\n", results)
}
