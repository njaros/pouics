package handler

import (
	"context"
	"pouic/internal/model"
)

type PlayerStore interface {
	GetById(ctx context.Context, id string) (*model.Player, error)
}

type PlayerHandler struct {
	players PlayerStore
}

func NewPlayerHandler(players PlayerStore) *PlayerHandler {
	return &PlayerHandler{
		players: players,
	}
}
