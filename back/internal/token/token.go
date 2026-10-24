package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalidToken is sent when a token is absent, expired, or not correctly signed.
var ErrInvalidToken = errors.New("invalid token")

// Manager creates and controls signed JWT with an HMAC secret.
type Manager struct {
	secret []byte
	ttl    time.Duration
}

// NewManager build a Manager with a given secret.
func NewManager(secret string) *Manager {
	return &Manager{
		secret: []byte(secret),
		ttl:    24 * time.Hour,
	}
}

// Generate create a signed JWT with the playerId as Subject.
func (m *Manager) Generate(playerId string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   playerId,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(m.secret)
}

// Verify controls the token validity and returns the playerId from it.
func (m *Manager) Verify(tokenString string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	t, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})
	if err != nil || !t.Valid {
		return "", ErrInvalidToken
	}
	return claims.Subject, nil
}
