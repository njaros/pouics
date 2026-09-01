package handler

import "pouic/internal/model"

type PlayerStore interface {
	GetById(id int) (*model.Player, error)
}

type PlayerHandler struct {
	players PlayerStore
}

func NewPlayerHandler(players PlayerStore) *PlayerHandler {
	return &PlayerHandler{
		players: players,
	}
}
