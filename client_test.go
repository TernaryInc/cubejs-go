package cube

// import (
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// )

// const (
// 	arbitraryDuration = time.Second * 1337
// 	cubeURL           = "cube/url"
// 	tenantID          = "Bruce_Tenant"
// 	token             = "super_secret_token"
// )

// type TestTokenGenerator struct {
// 	cubeAccessToken CubeAccessToken
// }

// func (ttg TestTokenGenerator) GenerateToken(tenantID string) (CubeAccessToken, error) {
// 	return ttg.cubeAccessToken, nil
// }

// type ErrorTokenGenerator struct {
// 	Error error
// }

// func (etg ErrorTokenGenerator) GenerateToken(tenantID string) (CubeAccessToken, error) {
// 	return CubeAccessToken{}, etg.Error
// }

// func Test_TokenNotSet(t *testing.T) {
// 	timeZero := time.Time{}
// 	expectedCubeAccessToken := CubeAccessToken{
// 		Token:     token,
// 		TenantID:  tenantID,
// 		ExpiresAt: timeZero,
// 	}
// 	testNower := util.TestNower{T: timeZero}
// 	tokenGenerator := TestTokenGenerator{expectedCubeAccessToken}

// 	cubeClient := NewCubeClient(tenantID, cubeURL, testNower, tokenGenerator)

// 	assert.Nil(t, cubeClient.accessToken)

// 	token, err := cubeClient.getToken()
// 	assert.Nil(t, err)
// 	assert.Equal(t, expectedCubeAccessToken.Token, token)
// 	assert.Equal(t, expectedCubeAccessToken.Token, cubeClient.accessToken.Token)
// }

// func Test_TokenExpiresSoon(t *testing.T) {
// 	var nower util.TestNower
// 	expiredCubeAccessToken := CubeAccessToken{
// 		Token:     token,
// 		TenantID:  tenantID,
// 		ExpiresAt: nower.Now().Add(timeToValidate - 1),
// 	}
// 	expectedCubeAccessToken := CubeAccessToken{
// 		Token:     token,
// 		TenantID:  tenantID,
// 		ExpiresAt: nower.Now().Add(timeToValidate + arbitraryDuration),
// 	}
// 	tokenGenerator := TestTokenGenerator{expectedCubeAccessToken}

// 	cubeClient := NewCubeClient(tenantID, cubeURL, nower, tokenGenerator)
// 	cubeClient.accessToken = &expiredCubeAccessToken

// 	assert.Equal(t, &expiredCubeAccessToken, cubeClient.accessToken)

// 	actualToken, err := cubeClient.getToken()
// 	assert.Nil(t, err)
// 	assert.Equal(t, expectedCubeAccessToken.Token, actualToken)
// 	assert.Equal(t, expectedCubeAccessToken.Token, cubeClient.accessToken.Token)
// }

// func Test_TokenNotExpired(t *testing.T) {
// 	var nower util.TestNower
// 	expectedCubeAccessToken := CubeAccessToken{
// 		Token:     token,
// 		TenantID:  tenantID,
// 		ExpiresAt: nower.Now().Add(timeToValidate + arbitraryDuration),
// 	}
// 	tokenGenerator := TestTokenGenerator{CubeAccessToken{}}

// 	cubeClient := NewCubeClient(tenantID, cubeURL, nower, tokenGenerator)
// 	cubeClient.accessToken = &expectedCubeAccessToken

// 	assert.Equal(t, &expectedCubeAccessToken, cubeClient.accessToken)

// 	actualToken, err := cubeClient.getToken()
// 	assert.Nil(t, err)
// 	assert.Equal(t, expectedCubeAccessToken.Token, actualToken)
// 	assert.Equal(t, expectedCubeAccessToken.Token, cubeClient.accessToken.Token)
// }

// func Test_TokenGetterError(t *testing.T) {
// 	expectedError := errors.New("get token error")
// 	tokenGenerator := ErrorTokenGenerator{expectedError}

// 	cubeClient := NewCubeClient(tenantID, cubeURL, util.TestNower{}, tokenGenerator)

// 	_, actualError := cubeClient.getToken()

// 	assert.NotNil(t, actualError)
// 	assert.Equal(t, expectedError, actualError)
// }
