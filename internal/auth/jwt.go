package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const issuer = "battleship-server"

type JWTManager struct {
	secret string
	TTL    time.Duration
}

func NewJWTManager(secret string, ttl time.Duration) *JWTManager {
	return &JWTManager{
		secret: secret,
		TTL:    ttl,
	}
}

func (m *JWTManager) Sign(playerID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   playerID,
		Issuer:    issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.TTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secret))
}

func (m *JWTManager) Verify(tokenStr string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method/alg: %v", t.Header["alg"])
		}
		return []byte(m.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %s, error: %v", tokenStr, err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("token parsed, but it's invalid: %v", token)
	}

	return claims, nil
}
