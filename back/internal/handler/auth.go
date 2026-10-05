package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"pouic/internal/common"
	"pouic/internal/dto"
	"pouic/internal/model"
	"pouic/internal/repository"
)

// AuthStore handles database operations.
type AuthStore interface {
	Create(ctx context.Context, name, passwordHash string) (*model.Player, error)
	GetByNameFull(ctx context.Context, name string) (*model.PlayerFull, error)
}

// TokenGenerator creates tokens.
type TokenGenerator interface {
	Generate(userID string) (string, error)
}

// AuthHandler handles submit and connection.
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

// Register creates an account and returns a token.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterDto
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

	user, err := h.players.Create(r.Context(), req.Name, passwordHash)
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

// Login control les credentials and send a token.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	player, err := h.players.GetByNameFull(r.Context(), req.Name)
	if err != nil {
		if errors.Is(err, repository.ErrPlayerNotFound) {
			common.WriteError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		common.WriteError(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	if !matchPassword(req.Password, player.Password) {
		common.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	h.respondWithToken(w, http.StatusOK, &player.Player)
}

// respondWithToken generates a token and serializes auth response.
func (h *AuthHandler) respondWithToken(w http.ResponseWriter, status int, player *model.Player) {
	tok, err := h.tokens.Generate(player.Id)
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	common.WriteJSON(w, status, dto.AuthResponse{Token: tok, Player: *player})
}

// hashPassword hash a password
func hashPassword(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

// matchPassword controls if the password is a good one.
func matchPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
