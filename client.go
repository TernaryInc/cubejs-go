/*
Package cube implements a simple client for Cube.
*/
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

	jsontime "github.com/liamylian/jsontime/v2/v2"
)

// Client is a client for interacting with Cube.
// Clients can be reused instead of created as needed.
// The methods of Client are safe for concurrent use by multiple goroutines.
type Client struct {
	tokenGenerator AccessTokenGenerator
	cubeURL        url.URL
}

// AccessTokenGenerator defines the interface for specifying access tokens which should be included in authenticated requests to the Cube server.
type AccessTokenGenerator interface {
	Get(ctx context.Context) (string, error)
}

// AccessTokenGeneratorFunc defines a simple concrete type which can satisfy the AccessTokenGenerator interface.
// This type is only one of many options, and users may define their own implementations of the interface.
type AccessTokenGeneratorFunc func(ctx context.Context) (string, error)

// Get generates an AccessToken or a non-nil error if it cannot.
func (fn AccessTokenGeneratorFunc) Get(ctx context.Context) (string, error) {
	return fn(ctx)
}

// NewClient creates a new Cube client.
// The optional tokenGenerator can be used to include an API token with the Cube requests.
func NewClient(cubeURL url.URL, tokenGenerator AccessTokenGenerator) *Client {
	return &Client{
		cubeURL:        cubeURL,
		tokenGenerator: tokenGenerator,
	}
}

// Load fetches JSON-encoded data from the Cube API and stores the result in the value pointed to by `results`. If `results` is nil or not a pointer, Load returns an error.
// Load uses the decodings that json.Unmarshal uses, allocating maps, slices, and pointers as necessary.
func (c *Client) Load(ctx context.Context, query Query, results interface{}) (ResponseMetadata, error) {
	var beginTime = time.Now()
	var requestBody = requestBody{query}

	marshaledRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return ResponseMetadata{}, fmt.Errorf("marshal load body: %w", err)
	}

	var (
		attempt            = 0
		continueWaitString = "Continue wait"
		// Do not spam Cube server with requests that will likely take some time
		// TODO: Replace with exponential backoff rate limiter
		limiter = ratelimit.New(1, ratelimit.Per(time.Minute))
	)

	for {
		var response *http.Response
		var responseBody responseBody
		attempt++

		var url = c.cubeURL
		url.Path += cubeLoadPath
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), bytes.NewBuffer(marshaledRequestBody))
		if err != nil {
			return ResponseMetadata{}, fmt.Errorf("new request with context: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		if c.tokenGenerator != nil {
			if token, err := c.tokenGenerator.Get(ctx); err != nil {
				return ResponseMetadata{}, fmt.Errorf("generate token: %w", err)
			} else {
				req.Header.Set("Authorization", token)
			}
		}

		limiter.Take()

		// Default timeout for queries is 10 minutes per cubejs documentation
		// https://cube.dev/docs/config#queue-options
		var netClient = &http.Client{
			Timeout: time.Minute * 10,
		}
		response, err = netClient.Do(req)
		if err != nil {
			return ResponseMetadata{}, fmt.Errorf("do request: %w", err)
		}

		defer func(body io.ReadCloser) {
			_ = body.Close()
		}(response.Body)

		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return ResponseMetadata{}, fmt.Errorf("read response bytes: %w", err)
		}

		if response.StatusCode >= 400 {
			var preview = strings.TrimSpace(string(responseBytes))
			return ResponseMetadata{}, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, preview[:min(1024, len(preview))])
		}

		err = json.Unmarshal(responseBytes, &responseBody)
		if err != nil {
			var preview = strings.TrimSpace(string(responseBytes))
			return ResponseMetadata{}, fmt.Errorf("decode response json (%s): %w", preview[:min(1024, len(preview))], err)
		}

		currentTime := time.Now()

		if responseBody.Error == "" {
			// TODO: Move into init func
			jsontime.SetDefaultTimeFormat("2006-01-02T15:04:05.000", time.UTC)
			var jsonTime = jsontime.ConfigWithCustomTimeFormat
			if err = jsonTime.Unmarshal(responseBody.Data, results); err != nil {
				return responseBody.ResponseMetadata, fmt.Errorf("unmarshal load response data: %w", err)
			}

			return responseBody.ResponseMetadata, nil
		} else if responseBody.Error != continueWaitString {
			return responseBody.ResponseMetadata, fmt.Errorf("load query results: %s", responseBody.Error)
		} else if currentTime.Sub(beginTime) > maximumQueryDuration {
			return responseBody.ResponseMetadata, fmt.Errorf("maximum query duration (%+v) exceeded after %d attempts", currentTime.Sub(beginTime), attempt)
		}
	}
}

// min returns the minimum of the two input ints
func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
