package token

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	userIDCorrect      = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	userLoginIDCorrect = "test0000"
)

func TestCreateAccessToken(t *testing.T) {
	tokenInfo, err := CreateAccessToken(&AuthClaims{UserID: userIDCorrect, UserLoginID: userLoginIDCorrect})
	require.NoError(t, err, "Failed to create access token")

	validatedAccessToken, err := ValidateAccessToken(tokenInfo.Token)
	require.NoError(t, err, "Failed to validate access token")
	require.Equal(t, validatedAccessToken.UserID, userIDCorrect)
	require.Equal(t, validatedAccessToken.UserLoginID, userLoginIDCorrect)
}

func TestCreateRefreshToken(t *testing.T) {
	tokenInfo, err := CreateRefreshToken(&AuthClaims{UserID: userIDCorrect, UserLoginID: userLoginIDCorrect})
	require.NoError(t, err, "Failed to create refresh token")

	validatedAccessToken, err := ValidateRefreshToken(tokenInfo.Token)
	require.NoError(t, err, "Failed to validate refresh token")
	require.Equal(t, validatedAccessToken.UserID, userIDCorrect)
	require.Equal(t, validatedAccessToken.UserLoginID, userLoginIDCorrect)
}
