package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sgallaghe1541/epilogue/internal/db"
)

const JWTExpired = "current token has expired"

type UserData struct {
}
type EpilogueClaims struct {
	Permissions []db.Permission `json:"permissions"`
	jwt.RegisteredClaims
}

func GenerateJWT(u *db.User, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)
	claims := &EpilogueClaims{
		u.Permissions,
		jwt.RegisteredClaims{
			Issuer:    "epilogue",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject:   fmt.Sprintf("%d", u.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret string) (string, []string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&EpilogueClaims{},
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return "", nil, err
	}

	claims, ok := token.Claims.(*EpilogueClaims)
	if !ok {
		return "", nil, errors.New("failed to parse claims")
	}

	expires, err := claims.GetExpirationTime()
	if err != nil {
		return "", nil, err
	}

	if expires.Time.Before(time.Now()) {
		return "", nil, errors.New(JWTExpired)
	}
	userIDString, err := claims.GetSubject()
	if err != nil {
		return "", nil, err
	}

	issuer, err := claims.GetIssuer()
	if err != nil {
		return "", nil, err
	}

	if issuer != string("epilogue") {
		return "", nil, errors.New("invalid issuer")
	}

	permissions := make([]string, len(claims.Permissions))
	for i, p := range claims.Permissions {
		permissions[i] = p.Stringify()
	}

	return userIDString, permissions, nil
}

func GenerateRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
