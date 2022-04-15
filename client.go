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
// TODO: Figure out an abstraction better than an interface
// Load fetches JSON-encoded data and stores the result in the value pointed to by `results`. If `results` is nil or not a pointer, Load returns an error. func (c *CubeClient) Load(ctx context.Context, query CubeQuery, results interface{}) (interface{}, error) {
// Load uses the decodings that json.Unmarshal uses, allocating maps, slices, and pointers as necessary.
func (c *CubeClient) Load(ctx context.Context, query CubeQuery, results interface{}) error {
	if err := query.Validate(); err != nil {
		return fmt.Errorf("invalid Cube query: %w", err)
	}

	var beginTime = time.Now()
	var loadBody = loadBody{query}

	var token = c.accessToken

	marshaledLoadBody, err := json.Marshal(loadBody)
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
		var loadResponse LoadResponse
		var response *http.Response
		attempt++

		var u = c.cubeURL
		u.Path = cubeLoadPath
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(marshaledLoadBody))
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

		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("read response bytes: %w", err)
		}

		if response.StatusCode >= 400 {
			return fmt.Errorf("unexpected status code %d: %s", response.StatusCode, strings.TrimSpace(string(bodyBytes)))
		}

		// TODO: Rename to request body and response body
		err = json.Unmarshal(bodyBytes, &loadResponse)
		if err != nil {
			return fmt.Errorf("decode response json (%s): %w", string(bodyBytes), err)
		}

		currentTime := time.Now()

		if loadResponse.Error == "" {
			// TODO: unmarshal loadResponse in the results pointer
			if err = json.Unmarshal(loadResponse.Data, results); err != nil {
				return fmt.Errorf("unmarshal load response data: %w", err)
			}

			return nil
		} else if loadResponse.Error != continueWaitString {
			return fmt.Errorf("load query results: %s", loadResponse.Error)
		} else if currentTime.Sub(beginTime) > maximumQueryDuration {
			return fmt.Errorf("maximum query duration (%+v) exceeded after %d attempts", currentTime.Sub(beginTime), attempt)
		}
	}
}
