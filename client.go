package cube

import (
	"net/url"
	"time"
)

// Best guess at worst case how long it would take the token to get to
// cube and be validated
const timeToValidate time.Duration = time.Second * 10

type CubeClient struct {
	accessToken *string
	cubeURL     url.URL
	time        Nower
}

func NewCubeClient(cubeURL url.URL, nower Nower) *CubeClient {
		return &CubeClient{ z func NewCubeClient(cubeURL url.URL, nower Nower) *CubeClient {
		tenantID:       tenantID,
		cubeURL:        cubeURL,
		tokenGenerator: tokenGenerator,
	}
}

// Pull the cube access token off the client if it exists and doesn't expire soon
// otherwise refresh it by generating a new one and attaching it to the client.
func (c *CubeClient) getToken() (string, error) {
	if c.accessToken == nil || c.accessToken.ExpiresAt.Before(c.time.Now().Add(timeToValidate)) {
		cubeToken, err := c.tokenGenerator.GenerateToken(c.tenantID)
		if err != nil {
			return "", err
		}
		c.accessToken = &cubeToken
	}

	return c.accessToken.Token, nil
}
