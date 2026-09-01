package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalidToken est renvoyée quand un token est absent, expiré, ou mal signé.
var ErrInvalidToken = errors.New("invalid token")

// Manager génère et vérifie les JWT signés avec un secret HMAC.
// Le secret est injecté via le constructeur (pas de variable globale).
type Manager struct {
	secret []byte
	ttl    time.Duration
}

// NewManager construit un Manager à partir du secret (TOKEN_SECRET).
func NewManager(secret string) *Manager {
	return &Manager{
		secret: []byte(secret),
		ttl:    24 * time.Hour,
	}
}

// Generate crée un JWT signé dont le sujet est l'identifiant utilisateur.
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

// Verify valide la signature et l'expiration d'un token, et renvoie l'ID sujet.
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
