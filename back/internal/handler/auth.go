package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"pouic/internal/common"
	"pouic/internal/model"
	"pouic/internal/repository"
)

// AuthStore décrit les opérations de persistance dont l'auth a besoin.
// L'interface est déclarée côté consommateur pour découpler le handler du
// repository concret (facilite les tests).
type AuthStore interface {
	Create(ctx context.Context, email, username, passwordHash string) (*model.Player, error)
	GetByName(ctx context.Context, email string) (*model.Player, error)
}

// TokenGenerator produit un token d'authentification pour un utilisateur.
type TokenGenerator interface {
	Generate(userID string) (string, error)
}

// AuthHandler gère l'inscription et la connexion.
type AuthHandler struct {
	players  AuthStore
	tokens   TokenGenerator
	validate *validator.Validate
}

// NewAuthHandler injecte les dépendances (store + générateur de token).
func NewAuthHandler(players AuthStore, tokens TokenGenerator) *AuthHandler {
	return &AuthHandler{
		players:  players,
		tokens:   tokens,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

// Register crée un compte, hashe le mot de passe et renvoie un token.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user, err := h.users.Create(r.Context(), req.Email, req.Username, passwordHash)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			common.WriteError(w, http.StatusConflict, "email or username already in use")
			return
		}
		common.WriteError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	h.respondWithToken(w, http.StatusCreated, user)
}

// Login vérifie les identifiants et renvoie un token.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.users.GetByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			// Même réponse que mot de passe invalide : on ne révèle pas
			// l'existence d'un compte.
			common.WriteError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		common.WriteError(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	if !matchPassword(req.Password, user.Password) {
		common.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	h.respondWithToken(w, http.StatusOK, user)
}

// respondWithToken génère un token et sérialise la réponse d'auth.
func (h *AuthHandler) respondWithToken(w http.ResponseWriter, status int, user *model.User) {
	tok, err := h.tokens.Generate(user.ID)
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	common.WriteJSON(w, status, model.AuthResponse{Token: tok, User: *user})
}

// hashPassword renvoie le hash bcrypt d'un mot de passe.
// bcrypt intègre un sel aléatoire ; le hash fait 60 caractères.
func hashPassword(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

// matchPassword vérifie qu'un mot de passe en clair correspond au hash stocké.
func matchPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
