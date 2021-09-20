package token

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	userIDCorrect      = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	userLoginIDCorrect = "test0000"
)

func TestCreateDeleteAuthToken(t *testing.T) {
	tokenInfo, err := CreateToken(&AuthClaims{UserID: userIDCorrect, UserLoginID: userLoginIDCorrect})
	require.NoError(t, err, "Failed to create auth token")

	validatedToken, err := ValidateToken(tokenInfo.Token)
	require.NoError(t, err, "Failed to validate auth token")
	require.Equal(t, validatedToken.UserID, userIDCorrect)
	require.Equal(t, validatedToken.UserLoginID, userLoginIDCorrect)
}
