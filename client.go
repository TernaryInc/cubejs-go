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

func NewCubeClient(cubeURL url.URL, accessToken *string) *CubeClient {
	return &CubeClient{
		cubeURL:     cubeURL,
		accessToken: accessToken,
	}
}

// TODO: Document
// Load fetches JSON-encoded data and stores the result in the value pointed to by `results`. If `results` is nil or not a pointer, Load returns an error. func (c *CubeClient) Load(ctx context.Context, query CubeQuery, results interface{}) (interface{}, error) {
// Load uses the decodings that json.Unmarshal uses, allocating maps, slices, and pointers as necessary.
func (c *CubeClient) Load(ctx context.Context, query CubeQuery, results interface{}) error {
	if err := query.Validate(); err != nil {
		return fmt.Errorf("invalid Cube query: %w", err)
	}

	var beginTime = time.Now()
	var requestBody = requestBody{query}
	var token = c.accessToken

	marshaledRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("marshal load body: %w", err)
	}

	var (
		attempt            = 0
		continueWaitString = "Continue wait"
		// Do not spam Cube server with requests that will likely take some time
		// TODO(bruce): Replace with exponential backoff rate limiter
		limiter = ratelimit.New(1, ratelimit.Per(time.Minute))
	)

	for {
		var response *http.Response
		var responseBody ResponseBody
		attempt++

		var u = c.cubeURL
		u.Path = cubeLoadPath
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(marshaledRequestBody))
		if err != nil {
			return fmt.Errorf("new request with context: %w", err)
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
			return fmt.Errorf("do request: %w", err)
		}

		defer func(body io.ReadCloser) {
			_ = body.Close()
		}(response.Body)

		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("read response bytes: %w", err)
		}

		if response.StatusCode >= 400 {
			var preview = strings.TrimSpace(string(responseBytes))
			return fmt.Errorf("unexpected status code %d: %s", response.StatusCode, preview[:min(1024, len(preview))])
		}

		err = json.Unmarshal(responseBytes, &responseBody)
		if err != nil {
			var preview = strings.TrimSpace(string(responseBytes))
			return fmt.Errorf("decode response json (%s): %w", preview[:min(1024, len(preview))], err)
		}

		currentTime := time.Now()

		if responseBody.Error == "" {
			// TODO: unmarshal loadResponse in the results pointer
			if err = json.Unmarshal(responseBody.Data, results); err != nil {
				return fmt.Errorf("unmarshal load response data: %w", err)
			}

			return nil
		} else if responseBody.Error != continueWaitString {
			return fmt.Errorf("load query results: %s", responseBody.Error)
		} else if currentTime.Sub(beginTime) > maximumQueryDuration {
			return fmt.Errorf("maximum query duration (%+v) exceeded after %d attempts", currentTime.Sub(beginTime), attempt)
		}
	}
}
