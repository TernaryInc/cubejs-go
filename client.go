package cube

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/ratelimit"
)

type CubeClient struct {
	// TODO: Swap out with anonymous function that can return a string or error
	accessToken *string
	cubeURL     url.URL
}

func NewCubeClient(cubeURL url.URL) *CubeClient {
	return &CubeClient{
		cubeURL: cubeURL,
	}
}

// TODO: Document
// TODO: Figure out an abstraction better than an interface
func (c *CubeClient) Load(ctx context.Context, query CubeQuery, results interface{}) (interface{}, error) {
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("invalid Cube query: %w", err)
	}

	var beginTime = time.Now()
	var loadBody = loadBody{query}

	var token = c.accessToken

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

	for {
		var loadResponse LoadResponse
		var response *http.Response
		attempt++

		var u = c.cubeURL
		u.Path = cubeLoadPath
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(marshaledLoadBody))
		if err != nil {
			return nil, fmt.Errorf("new request with context: %w", err)
		}

		// Allow the user to leave off a token
		// TODO: test
		if token != nil {
			req.Header.Set("Authorization", *token)
		}
		req.Header.Set("Content-Type", "application/json")

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

		err = json.Unmarshal(bodyBytes, &loadBody)
		if err != nil {
			return nil, fmt.Errorf("decode response json (%s): %w", string(bodyBytes), err)
		}

		currentTime := time.Now()

		if loadResponse.Error == "" {
			return &loadResponse, nil
		} else if loadResponse.Error != continueWaitString {
			return nil, fmt.Errorf("load query results: %s", loadResponse.Error)
		} else if currentTime.Sub(beginTime) > maximumQueryDuration {
			return nil, fmt.Errorf("maximum query duration (%+v) exceeded after %d attempts", currentTime.Sub(beginTime), attempt)
		}
	}
}
