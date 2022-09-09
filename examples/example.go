package examples

import (
	"context"
	"fmt"
	"net/url"
	"time"

	cube "github.com/TernaryInc/cubejs-go"
)

const exampleAccessToken = "This is a fake example access token"

type tokenGenerator struct{}

func (generator tokenGenerator) Get(ctx context.Context) (string, error) {
	return exampleAccessToken, nil
}

// TestLocalServer demonstrates how to initialize a Cube client and use it to query the Cube service, binding the results to a slice of Go structs.
func TestLocalServer() error {
	var u = url.URL{
		Scheme: "http",
		Host:   "localhost:4000",
	}

	var tokenGenerator = tokenGenerator{}

	var cubeClient = cube.NewClient(u, tokenGenerator)

	var cubeQuery = cube.Query{
		Measures:   []string{"MyCube.measure1"},
		Dimensions: []string{"MyCube.dimension1"},
		TimeDimensions: []cube.TimeDimension{
			{
				Dimension: "MyCube.timestamp",
				DateRange: cube.DateRange{
					AbsoluteRange: []string{
						"2022-04-01",
						"2022-04-20",
					},
				},
				Granularity: cube.Granularity_Day,
			},
		},
	}

	type QueryResult struct {
		Cost      float64   `json:"MyCube.measure1"`
		ProjectID string    `json:"MyCube.dimension1"`
		Timestamp time.Time `json:"MyCube.timestamp"`
	}

	var results []QueryResult

	var _, err = cubeClient.Load(context.Background(), cubeQuery, &results)
	if err != nil {
		return fmt.Errorf("load query: %w", err)
	}

	return nil
}
