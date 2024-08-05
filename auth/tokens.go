package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const JWTExpired = "current token has expired"

type EpilogueClaims struct {
	Div string `json:"div"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int, div, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)
	claims := EpilogueClaims{
		div,
		jwt.RegisteredClaims{
			Issuer:    "epilogue",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret string) (string, string, error) {
	claimsStruct := EpilogueClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return "", "", err
	}

	expires, err := token.Claims.GetExpirationTime()
	if err != nil {
		return "", "", err
	}

	if expires.Time.Before(time.Now()) {
		return "", "", errors.New(JWTExpired)
	}
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", "", err
	}

	userDiv := claimsStruct.Div
	if userDiv == "" {
		return "", "", errors.New("division didn't come through")
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", "", err
	}

	if issuer != string("epilogue") {
		return "", "", errors.New("invalid issuer")
	}

	return userIDString, userDiv, nil
}

func GenerateRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
