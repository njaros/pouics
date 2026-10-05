package player

import (
	"context"
	"errors"
	"pouic/internal/model"
	"pouic/internal/repository"
	"pouic/internal/repository/common"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pgHandler *common.PgHandler
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pgHandler: common.NewPgHandler(pool)}
}

func (r *Repository) Create(ctx context.Context, name string, passwordHashed string) (*model.Player, error) {
	var player model.Player
	err := r.pgHandler.Pool.QueryRow(ctx, queryCreate, name, passwordHashed).Scan(&player.Id, &player.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repository.ErrPlayerNotFound
	}
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *Repository) GetById(ctx context.Context, id string) (*model.Player, error) {
	var player model.Player
	err := r.pgHandler.Pool.QueryRow(ctx, queryGetById, id).Scan(&player.Id, &player.Name)
	return r.pgHandler.HandleGetResponse(&player, err)
}

func (r *Repository) GetByName(ctx context.Context, name string) (*model.Player, error) {
	var player model.Player
	err := r.pgHandler.Pool.QueryRow(ctx, queryGetByName, name).Scan(&player.Id, &player.Name)
	return r.pgHandler.HandleGetResponse(&player, err)
}

func (r *Repository) GetByNameFull(ctx context.Context, name string) (*model.PlayerFull, error) {
	var player model.PlayerFull
	err := r.pgHandler.Pool.QueryRow(ctx, queryGetByNameFull, name).Scan(&player.Id, &player.Name, &player.Password)
	return r.pgHandler.HandleGetResponse(&player, err)
}