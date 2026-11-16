package game

import (
	"context"
	"log"
	"pouic/internal/websocket/command"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Hub interface {
	GetCommands() chan command.Command
	EmitTo(playerId string, data interface{})
}

type Engine struct {
	s State
	stateLocker sync.Mutex
	pool *pgxpool.Pool
	hub Hub
	lastSave time.Time
}

func NewEngine(pool *pgxpool.Pool, hub Hub) *Engine {
	return &Engine{
		s: *NewState(pool),
		stateLocker: sync.Mutex{},
		pool: pool,
		hub: hub,
		lastSave: time.Now(),
	}
}

// Load the engine state from db, then run the engine goroutine.
func (e *Engine) Run() {
	e.s.LoadFromDb(context.Background());
	go e.run();
}

// Run the loop of the sim engine.
// It start the Command Reader in a goroutine.
// It calls the engine routine in a loop, and wait for at least one second to do the next loop.
// If the routine is longer than one second, that is not cool, and the commands will not be read entirely.
func (e *Engine) run() {
	timer := time.NewTimer(time.Second);
	go e.commandsReaderGoroutine();
	for {
		chrono := time.Now();
		e.routine();
		ellapsed := time.Since(chrono);
		log.Println("routine time ellapsed: %d milliseconds.", ellapsed.Milliseconds());
		select {
		case <- timer.C:
			timer.Reset(time.Second);
		}
	}
}

func (e *Engine) commandsReaderGoroutine() {
	ticker := time.NewTicker(20 * time.Millisecond);
	for {
		e.stateLocker.Lock();
		select {
		case <-e.hub.GetCommands():
			//TODO commandparsing
			e.stateLocker.Unlock();
		case <-ticker.C:
			e.stateLocker.Unlock();
		}
	}
}

// Listen into query channels
// Update the states
// Broadcast / emit things
// The routine admits that the State is static and he can do whatever he wants to it,
// So the Command Reader must be stopped when the routine is running.
func (e *Engine) routine() {
	defer e.stateLocker.Unlock()
	colonies := e.s.GetColonies()
	for _, colony := range colonies {
		colony.Move()
	}
	e.stateLocker.Lock()
}
