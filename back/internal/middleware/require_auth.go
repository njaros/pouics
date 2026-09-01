package middleware

import (
	"context"
	"pouic/internal/common"
	"net/http"
	"strings"
)

// TokenVerifier valide un token et renvoie l'identifiant de l'utilisateur.
// L'interface est déclarée côté consommateur pour découpler le middleware du
// token.Manager concret (facilite les tests).
type TokenVerifier interface {
	Verify(tokenString string) (string, error)
}

// contextKey est un type privé pour éviter les collisions de clés de contexte.
type contextKey string

const playerIdKey contextKey = "playerId"

// Authenticator porte la dépendance de vérification des tokens.
type Authenticator struct {
	tokens TokenVerifier
}

// NewAuthenticator injecte le vérificateur de token (pas de variable globale).
func NewAuthenticator(tokens TokenVerifier) *Authenticator {
	return &Authenticator{tokens: tokens}
}

// RequireAuth est un middleware chi qui exige un token Bearer valide.
// Il extrait le token du header Authorization, le vérifie, puis injecte
// l'identifiant utilisateur dans le contexte de la requête. En cas d'échec,
// il répond 401 et n'appelle pas le handler suivant.
func (a *Authenticator) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := bearerToken(r)
		if !ok {
			common.Unauthorized(w)
			return
		}

		playerId, err := a.tokens.Verify(tokenString)
		if err != nil {
			common.Unauthorized(w)
			return
		}

		ctx := context.WithValue(r.Context(), playerIdKey, playerId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserIDFromContext récupère l'identifiant utilisateur injecté par RequireAuth.
// Le booléen est faux si la requête n'est pas passée par le middleware d'auth.
func PlayerIdFromContext(ctx context.Context) (string, bool) {
	playerId, ok := ctx.Value(playerIdKey).(string)
	return playerId, ok
}

// bearerToken extrait le token du header "Authorization: Bearer <token>".
func bearerToken(r *http.Request) (string, bool) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", false
	}
	const prefix = "Bearer "
	if len(header) <= len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return "", false
	}
	return strings.TrimSpace(header[len(prefix):]), true
}