package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/model"
)

const (
	accTokenKey        = "cyQsIQ6RSE1CqTARl8pWeM7br9qp1Don57Pd18uDCwoBaiUPEXWe15pYMP4D9WKc"
	accTokenTimeoutMin = 60 // 1 hour
	refTokenKey        = "lEAWYT9pcR5r9B5fq3ED2V5dQyhZlOACZD0lJJwzMmzxScOAX1k1ZuXHZ9hLAOG9"
	refTokenTimeoutMin = 60 * 24 * 14 // 2 weeks
)

// Structs
type TokenClaims struct {
	jwt.StandardClaims
	AuthClaims
}

type AuthClaims struct {
	UserID      string
	UserLoginID string
	UserRole    model.UserRole
}

type TokenInfo struct {
	Token     string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func CreateAccessToken(authInfo *AuthClaims) (*TokenInfo, error) {
	return createToken(accTokenKey, accTokenTimeoutMin, authInfo)
}

func CreateRefreshToken(authInfo *AuthClaims) (*TokenInfo, error) {
	return createToken(refTokenKey, refTokenTimeoutMin, authInfo)
}

func createToken(tokenKey string, tokenTimeoutMin int, authInfo *AuthClaims) (*TokenInfo, error) {
	// Calculate issuance and expiration time
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(time.Minute * accTokenTimeoutMin)

	// Set access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  issuedAt.Unix(),
			ExpiresAt: expiresAt.Unix(),
		},
		AuthClaims: AuthClaims{
			UserID:      authInfo.UserID,
			UserLoginID: authInfo.UserLoginID,
			UserRole:    authInfo.UserRole,
		},
	})

	// Signing access token
	tokenSigned, err := token.SignedString([]byte(tokenKey))
	if err != nil {
		return nil, err
	}

	return &TokenInfo{
		Token:     tokenSigned,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
	}, nil
}

func ValidateAccessToken(token string) (*AuthClaims, error) {
	return validateToken(accTokenKey, token)
}

func ValidateRefreshToken(token string) (*AuthClaims, error) {
	return validateToken(refTokenKey, token)
}

func validateToken(tokenKey, tokenSigned string) (*AuthClaims, error) {
	// Prase token
	claims := TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenSigned, &claims, func(token *jwt.Token) (interface{}, error) {
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
