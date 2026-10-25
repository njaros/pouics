package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"pouic/internal/handler"
	myMiddlewares "pouic/internal/middleware"
	"pouic/internal/repository/player"
	"pouic/internal/token"
	"pouic/internal/websocket/command"
	"pouic/internal/websocket/hub"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Apply the migrations in the migrations folder.
// TODO: a better system to force any migration.
func db_migration(database_url string) {
	force := os.Getenv("MIGRATION_FORCE")
	m, err := migrate.New("file://migrations", database_url)
	if err != nil {
		log.Fatal("failed to migrate: ", err)
	}
	if force == "true" {
		to_str := os.Getenv("MIGRATION_FORCE_TO")
		if to, err := strconv.Atoi(to_str); err != nil {
			log.Fatal("unable to parse the migrations version: ", to_str)
		} else if err := m.Force(to); err != nil {
			log.Fatal("failed to m.Force(): ", err)
		}
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to m.Up(): ", err)
	}
}

// Get a pointer to the db pool to send concurrency request to the db.
func get_db_pool(ctx context.Context) *pgxpool.Pool {
	database_url := os.Getenv("DATABASE_URL")
	if database_url == "" {
		log.Fatal("DATABASE_URL env variable not set.")
	}
	db_migration(database_url)
	pool, err := pgxpool.New(ctx, database_url)
	if err != nil {
		log.Fatal("failed to pgxpool.New", err)
	}
	return pool
}

func main() {
	pool := get_db_pool(context.Background())
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		AllowCredentials: true,
		MaxAge: 300,
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/health"))
	r.Use(middleware.Recoverer)

	token_secret := os.Getenv("TOKEN_SECRET")
	if token_secret == "" {
		log.Fatal("TOKEN_SECRET is not set")
	}

	// Token management set up
	tokenManager := token.NewManager(token_secret)
	validator := myMiddlewares.NewAuthenticator(tokenManager)

	// Websockets set up
	parser := command.NewParser(make(map[string]any))
	hub := hub.NewHub()

	// API routes set up
	authStore := player.NewRepository(pool)
	authHandler := handler.NewAuthHandler(authStore, tokenManager)
	// playerHandler := handler.NewPlayerHandler(authStore)

	r.Group(func(r chi.Router) {
		r.Use(validator.RequireAuth)
		r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
			handler.ServeWs(hub, parser, w, r)
		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello"))
		})
	})

	r.Get("/doc", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"body": "ceci est la doc"}`))
	})

	r.Group(func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	http.ListenAndServe(":3333", r)
}