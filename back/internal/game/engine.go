package game

import (
	"pouic/internal/websocket/command"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Data interface {}

type Hub interface {
	GetCommands() chan command.Command
	EmitTo(playerId string, data Data)
}

type Engine struct {
	s State
	pool *pgxpool.Pool
	hub Hub
	lastSave time.Time
}