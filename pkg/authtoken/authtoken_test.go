package authtoken

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	userIDCorrect      = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	userLoginIDCorrect = "test0000"
)

func TestCreateDeleteAuthToken(t *testing.T) {
	token, err := CreateAuthToken(&AuthInfo{UserID: userIDCorrect, UserLoginID: userLoginIDCorrect})
	require.NoError(t, err, "Failed to create auth token")

	tokenInfo, err := ValidateAuthToken(token)
	require.NoError(t, err, "Failed to validate auth token")
	require.Equal(t, tokenInfo.UserID, userIDCorrect)
	require.Equal(t, tokenInfo.UserLoginID, userLoginIDCorrect)
}
