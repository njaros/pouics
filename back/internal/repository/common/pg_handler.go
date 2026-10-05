package common

import (
	"errors"
	"pouic/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgHandler struct {
	Pool *pgxpool.Pool
}

func NewPgHandler(pool *pgxpool.Pool) *PgHandler {
	return &PgHandler{Pool: pool}
}

func (p *PgHandler) HandleGetResponse[T any](obj *T, err error) (*T, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrPlayerNotFound
	}
	if err != nil {
		return nil, err
	}
	return obj, nil
}