package game

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Engine struct {
	s State
	pool *pgxpool.Pool
	chann chan PlayerRequest
	lastSave time.Time
}