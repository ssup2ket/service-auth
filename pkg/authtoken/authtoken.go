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

type AuthClaim struct {
	jwt.StandardClaims
	AuthInfo
}

type AuthInfo struct {
	UserID      string
	UserLoginID string
}

func CreateAuthToken(authInfo *AuthInfo) (string, error) {
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &AuthClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		},
		AuthInfo: AuthInfo{
			UserID:      authInfo.UserID,
			UserLoginID: authInfo.UserLoginID,
		},
	})

	// Signing token
	signedToken, err := token.SignedString([]byte(tokenKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateAuthToken(signedToken string) (*AuthInfo, error) {
	// Prase token
	claims := AuthClaim{}
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
	return &claims.AuthInfo, nil
}
