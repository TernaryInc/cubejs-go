package cube_test

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

const cubeTokenTTL = time.Hour

type CubeAccessToken struct {
	Token     string    `json:"token"`
	TenantID  string    `json:"tenantID"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type CubeTokenGenerator struct {
	cubeAPISecret string
}

type TokenGenerator interface {
	GenerateToken(string) (CubeAccessToken, error)
}

type cubeUserContext struct {
	TenantID string `json:"tenantID"`
}

type cubeClaims struct {
	cubeUserContext
	jwt.StandardClaims
}

func NewTokenGenerator(cubeAPISecret string) (CubeTokenGenerator, error) {
	if cubeAPISecret == "" {
		return CubeTokenGenerator{}, errors.New("cubeAPISecret cannot be empty")
	}

	return CubeTokenGenerator{cubeAPISecret: cubeAPISecret}, nil
}

// GenerateToken issues a cube token scoped to the passed in tenant ID. Assume all security checks
// have already passed.
func (ctg CubeTokenGenerator) GenerateToken(tenantID string) (CubeAccessToken, error) {
	var expiration = time.Now().Add(cubeTokenTTL)
	var tokenMaker = jwt.NewWithClaims(jwt.SigningMethodHS256, cubeClaims{
		// This claim is read by Cube.JS, and entitles the bearer to any BigQuery dataset in data_${tenantID}
		cubeUserContext: cubeUserContext{
			TenantID: tenantID,
		},

		// Standard claim for expiration
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	})

	var token, err = tokenMaker.SignedString([]byte(ctg.cubeAPISecret))
	if err != nil {
		return CubeAccessToken{}, err
	}

	var cubeToken = CubeAccessToken{
		ExpiresAt: expiration,
		Token:     token,
		TenantID:  tenantID,
	}

	return cubeToken, nil
}
