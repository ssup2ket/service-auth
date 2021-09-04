package authtoken

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	tokenKey        = "ssup2ket"
	tokenTimeoutMin = 60
)

// Structs
type AuthTokenClaims struct {
	jwt.StandardClaims
	AuthClaims
}

type AuthClaims struct {
	UserID      string
	UserLoginID string
}

type AuthTokenInfo struct {
	Token     string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func CreateAuthToken(authInfo *AuthClaims) (*AuthTokenInfo, error) {
	// Calcuation issuance and expiration time
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(time.Minute * 60)

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &AuthTokenClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  issuedAt.Unix(),
			ExpiresAt: expiresAt.Unix(),
		},
		AuthClaims: AuthClaims{
			UserID:      authInfo.UserID,
			UserLoginID: authInfo.UserLoginID,
		},
	})

	// Signing token
	signedToken, err := token.SignedString([]byte(tokenKey))
	if err != nil {
		return nil, err
	}
	return &AuthTokenInfo{
		Token:     signedToken,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}, nil
}

func ValidateAuthToken(signedToken string) (*AuthClaims, error) {
	// Prase token
	claims := AuthTokenClaims{}
	token, err := jwt.ParseWithClaims(signedToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey), nil
	})
	if err != nil {
		return nil, err
	}

	// Check token's validation
	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	// Return auth infos
	return &claims.AuthClaims, nil
}
