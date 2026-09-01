package cmd

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Apply the migrations in the migrations folder.
func db_migration(database_url string) {
	force := os.Getenv("MIGRATION_FORCE")
	m, err := migrate.New("file://migrations", database_url)
	if err != nil {
		log.Fatal("failed to migrate.New", err)
	}
	if force == "true" {
		to_str := os.Getenv("MIGRATION_FORCE_TO")
		if to, err := strconv.Atoi(to_str); err != nil {
			log.Fatal("unable to parse the migrations version: ", to_str)
		} else if err := m.Force(to); err != nil {
			log.Fatal("failed to m.Force()", err)
		}
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to m.Up()", err)
	}
}

// Get a pointer to the db pool to send concurrency request to the db.
func get_db_pool(ctx context.Context) *pgxpool.Pool {
	database_url := os.Getenv("DATABASE_URL")
	db_migration(database_url)
	pool, err := pgxpool.New(ctx, database_url)
	if err != nil {
		log.Fatal("failed to pgxpool.New", err)
	}
	return pool
}

func main() {
	r := chi.NewRouter()
	pool := get_db_pool(context.Background())

	token_secret := os.Getenv("TOKEN_SECRET")
	if token_secret == "" {
		log.Fatal("TOKEN_SECRET is not set")
	}

	r.Group()
}